// Package cli implementa a interface de linha de comando (spec 002) e o modo
// interativo (spec 005). Toda regra de negócio fica em internal/app.
package cli

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/fabiodrneles/cv-craft/internal/app"
	"github.com/fabiodrneles/cv-craft/internal/generator"
	"github.com/fabiodrneles/cv-craft/internal/i18n"
	"github.com/fabiodrneles/cv-craft/internal/resume"
)

// Exit codes (spec 002).
const (
	ExitOK       = 0
	ExitError    = 1
	ExitUsage    = 2
	ExitExists   = 3
	ExitInvalid  = 4
	defaultInput = "curriculum.yaml"
)

// Env é o ambiente do processo, injetável para testes.
type Env struct {
	Stdin          io.Reader
	Stdout, Stderr io.Writer
	// IsTerminal informa se stdin é um terminal interativo.
	IsTerminal func() bool
	Getenv     func(string) string

	Version, Commit, Date string
}

type runner struct {
	env Env
	in  *bufio.Reader
	out io.Writer
	err io.Writer
}

// Run executa a CLI com os argumentos (sem o nome do programa) e devolve o
// exit code.
func Run(args []string, env Env) int {
	if env.Stdin == nil {
		env.Stdin = strings.NewReader("")
	}
	if env.IsTerminal == nil {
		env.IsTerminal = func() bool { return false }
	}
	if env.Getenv == nil {
		env.Getenv = func(string) string { return "" }
	}
	if env.Version == "" {
		env.Version = "dev"
	}
	r := &runner{env: env, in: bufio.NewReader(env.Stdin), out: env.Stdout, err: env.Stderr}
	return r.dispatch(args)
}

func (r *runner) dispatch(args []string) int {
	if len(args) == 0 {
		fmt.Fprint(r.out, usage)
		return ExitOK
	}
	cmd, rest := args[0], args[1:]
	switch cmd {
	case "build":
		return r.build(rest)
	case "validate":
		return r.validate(rest)
	case "init":
		return r.init(rest)
	case "version", "--version", "-version":
		return r.version()
	case "help", "-h", "--help", "-help":
		return r.help(rest)
	case "ui", "interactive":
		return r.interactive()
	}
	r.errorf("comando desconhecido %q", cmd)
	fmt.Fprintln(r.err, "Use 'cv-craft help' para ver os comandos disponíveis.")
	return ExitUsage
}

func (r *runner) build(args []string) int {
	fs := newFlagSet("build")
	var format, output, lang string
	var force, quiet, verbose bool
	stringFlag(fs, &format, "pdf", "format", "f")
	stringFlag(fs, &output, "", "output", "o")
	stringFlag(fs, &lang, "", "lang")
	boolFlag(fs, &force, "force", "y")
	boolFlag(fs, &quiet, "quiet", "q")
	boolFlag(fs, &verbose, "verbose", "v")

	pos, code, ok := r.parse(fs, args, usageBuild)
	if !ok {
		return code
	}
	if len(pos) != 1 {
		return r.usageError(usageBuild, "informe exatamente um arquivo YAML")
	}

	var formats []generator.Format
	if strings.EqualFold(format, "all") {
		formats = generator.Formats
	} else {
		f, err := generator.ParseFormat(format)
		if err != nil {
			return r.usageError(usageBuild, err.Error())
		}
		formats = []generator.Format{f}
	}
	if lang != "" {
		if _, err := i18n.Get(lang); err != nil {
			return r.usageError(usageBuild, err.Error())
		}
	}

	opts := app.BuildOptions{
		Input:     pos[0],
		Output:    output,
		Formats:   formats,
		Lang:      lang,
		Force:     force,
		Version:   r.env.Version,
		CreatedAt: r.sourceDate(),
	}
	if r.env.IsTerminal() {
		opts.Confirm = r.confirm
	}
	if verbose {
		fmt.Fprintf(r.err, "entrada: %s | formatos: %s\n", opts.Input, joinFormats(formats))
	}

	res, err := app.Build(opts)
	if !quiet {
		r.warnings(res.Warnings)
	}
	if err != nil {
		return r.fail(err)
	}
	if verbose {
		fmt.Fprintf(r.err, "idioma: %s\n", res.Lang)
	}
	if !quiet {
		for _, f := range res.Files {
			fmt.Fprintf(r.out, "Gerado: %s (%s)\n", f.Path, humanSize(f.Size))
		}
	}
	return ExitOK
}

func (r *runner) validate(args []string) int {
	fs := newFlagSet("validate")
	pos, code, ok := r.parse(fs, args, usageValidate)
	if !ok {
		return code
	}
	if len(pos) != 1 {
		return r.usageError(usageValidate, "informe exatamente um arquivo YAML")
	}
	cv, rep, err := app.Load(pos[0])
	r.warnings(rep.Warnings)
	if err != nil {
		return r.fail(err)
	}
	fmt.Fprintf(r.out, "Válido: %s — %s | %d experiência(s), %d categoria(s) de habilidades, %d formação(ões)\n",
		pos[0], cv.Contact.Name, len(cv.Experience), len(cv.Skills), len(cv.Education))
	return ExitOK
}

func (r *runner) init(args []string) int {
	fs := newFlagSet("init")
	var force bool
	boolFlag(fs, &force, "force", "y")
	pos, code, ok := r.parse(fs, args, usageInit)
	if !ok {
		return code
	}
	if len(pos) > 1 {
		return r.usageError(usageInit, "informe no máximo um arquivo")
	}
	path := defaultInput
	if len(pos) == 1 {
		path = pos[0]
	}
	if err := app.Init(path, force); err != nil {
		return r.fail(err)
	}
	fmt.Fprintf(r.out, "Criado: %s\nEdite o arquivo e gere o currículo com: cv-craft build %s\n", path, quoteArg(path))
	return ExitOK
}

func (r *runner) version() int {
	fmt.Fprintf(r.out, "cv-craft %s", r.env.Version)
	if r.env.Commit != "" || r.env.Date != "" {
		fmt.Fprintf(r.out, " (commit %s, %s)", orDefault(r.env.Commit, "desconhecido"), orDefault(r.env.Date, "data desconhecida"))
	}
	fmt.Fprintln(r.out)
	return ExitOK
}

func (r *runner) help(args []string) int {
	if len(args) == 0 {
		fmt.Fprint(r.out, usage)
		return ExitOK
	}
	text, ok := commandUsage[args[0]]
	if !ok {
		return r.usageError(usage, fmt.Sprintf("comando desconhecido %q", args[0]))
	}
	fmt.Fprint(r.out, text)
	return ExitOK
}

// ---- helpers ----

func (r *runner) parse(fs *flag.FlagSet, args []string, help string) (pos []string, code int, ok bool) {
	pos, err := parseInterspersed(fs, args)
	if errors.Is(err, flag.ErrHelp) {
		fmt.Fprint(r.out, help)
		return nil, ExitOK, false
	}
	if err != nil {
		return nil, r.usageError(help, translateFlagError(err)), false
	}
	return pos, 0, true
}

func (r *runner) usageError(help, msg string) int {
	r.errorf("%s", msg)
	fmt.Fprintln(r.err)
	fmt.Fprint(r.err, help)
	return ExitUsage
}

func (r *runner) errorf(format string, args ...any) {
	fmt.Fprintf(r.err, "erro: "+format+"\n", args...)
}

func (r *runner) warnings(ws []resume.Issue) {
	for _, w := range ws {
		fmt.Fprintf(r.err, "aviso: %s\n", w)
	}
}

// fail converte um erro em mensagem e exit code.
func (r *runner) fail(err error) int {
	var inv *app.InvalidError
	switch {
	case errors.As(err, &inv):
		if inv.Err != nil {
			r.errorf("%v", inv)
		}
		if len(inv.Report.Errors) > 0 {
			r.errorf("%s tem %d problema(s):", inv.Path, len(inv.Report.Errors))
			for _, e := range inv.Report.Errors {
				fmt.Fprintf(r.err, "  - %s\n", e)
			}
		}
		return ExitInvalid
	case errors.Is(err, app.ErrCanceled):
		fmt.Fprintln(r.err, "Cancelado. Nenhum arquivo foi alterado.")
		return ExitExists
	case errors.Is(err, app.ErrOutputExists):
		r.errorf("%v", err)
		return ExitExists
	}
	r.errorf("%v", err)
	return ExitError
}

// confirm pergunta, no terminal, se um arquivo pode ser sobrescrito.
func (r *runner) confirm(path string) (bool, error) {
	fmt.Fprintf(r.err, "O arquivo %s já existe. Sobrescrever? [s/N] ", path)
	line, err := r.in.ReadString('\n')
	if err != nil && line == "" {
		fmt.Fprintln(r.err)
		return false, nil
	}
	switch strings.ToLower(strings.TrimSpace(line)) {
	case "s", "sim", "y", "yes":
		return true, nil
	}
	return false, nil
}

// sourceDate respeita SOURCE_DATE_EPOCH para builds reproduzíveis.
func (r *runner) sourceDate() time.Time {
	if s := r.env.Getenv("SOURCE_DATE_EPOCH"); s != "" {
		if n, err := strconv.ParseInt(s, 10, 64); err == nil {
			return time.Unix(n, 0).UTC()
		}
	}
	return time.Time{}
}

func newFlagSet(name string) *flag.FlagSet {
	fs := flag.NewFlagSet(name, flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	return fs
}

func stringFlag(fs *flag.FlagSet, p *string, def string, names ...string) {
	for _, n := range names {
		fs.StringVar(p, n, def, "")
	}
}

func boolFlag(fs *flag.FlagSet, p *bool, names ...string) {
	for _, n := range names {
		fs.BoolVar(p, n, false, "")
	}
}

// parseInterspersed permite flags antes e depois dos argumentos posicionais
// (spec 002, AC-6).
func parseInterspersed(fs *flag.FlagSet, args []string) ([]string, error) {
	var pos []string
	for {
		if err := fs.Parse(args); err != nil {
			return nil, err
		}
		rest := fs.Args()
		if len(rest) == 0 {
			return pos, nil
		}
		pos = append(pos, rest[0])
		args = rest[1:]
	}
}

func translateFlagError(err error) string {
	msg := err.Error()
	if f, ok := strings.CutPrefix(msg, "flag provided but not defined: "); ok {
		return "flag desconhecida: " + f
	}
	if f, ok := strings.CutPrefix(msg, "flag needs an argument: "); ok {
		return "a flag " + f + " precisa de um valor"
	}
	return msg
}

func joinFormats(fs []generator.Format) string {
	s := make([]string, len(fs))
	for i, f := range fs {
		s[i] = string(f)
	}
	return strings.Join(s, ", ")
}

func humanSize(n int64) string {
	if n < 1024 {
		return fmt.Sprintf("%d B", n)
	}
	return fmt.Sprintf("%.1f KB", float64(n)/1024)
}

func orDefault(s, def string) string {
	if s == "" {
		return def
	}
	return s
}

func quoteArg(s string) string {
	if strings.ContainsAny(s, " \t") {
		return `"` + s + `"`
	}
	return s
}
