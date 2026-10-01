// Package app concentra os casos de uso do CV-Craft (build, validate, init).
// A CLI e o modo interativo apenas traduzem argumentos para estas funções
// (spec 002, FR-2).
package app

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fabiodrneles/cv-craft/examples"
	"github.com/fabiodrneles/cv-craft/internal/generator"
	"github.com/fabiodrneles/cv-craft/internal/i18n"
	"github.com/fabiodrneles/cv-craft/internal/resume"
)

var (
	// ErrOutputExists indica que a saída já existe e não houve confirmação.
	ErrOutputExists = errors.New("arquivo já existe")
	// ErrCanceled indica que o usuário recusou sobrescrever.
	ErrCanceled = errors.New("operação cancelada")
)

// InvalidError indica um YAML que não pôde ser lido ou não passou na validação.
type InvalidError struct {
	Path   string
	Err    error         // erro de leitura/sintaxe; nil se só houve erros de validação
	Report resume.Report // erros e avisos de validação
}

func (e *InvalidError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Path, e.Err)
	}
	return fmt.Sprintf("%s: %d erro(s) de validação", e.Path, len(e.Report.Errors))
}

func (e *InvalidError) Unwrap() error { return e.Err }

// Load lê, decodifica e valida um currículo. Erros de leitura de arquivo são
// devolvidos como estão; problemas no conteúdo vêm como *InvalidError.
func Load(path string) (*resume.Resume, resume.Report, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, resume.Report{}, fmt.Errorf("arquivo não encontrado: %s", path)
		}
		return nil, resume.Report{}, err
	}
	r, rep, err := resume.Parse(data)
	if err != nil {
		return nil, rep, &InvalidError{Path: path, Err: err, Report: rep}
	}
	rep.Merge(resume.Validate(r))
	if !rep.OK() {
		return r, rep, &InvalidError{Path: path, Report: rep}
	}
	return r, rep, nil
}

// BuildOptions configura uma geração.
type BuildOptions struct {
	Input   string
	Output  string // arquivo (um formato) ou diretório (vários formatos); vazio = ao lado do YAML
	Formats []generator.Format
	Lang    string // sobrepõe meta.locale
	Force   bool
	// Confirm pergunta se um arquivo existente pode ser sobrescrito. Nil
	// significa ambiente não interativo: a geração falha com ErrOutputExists.
	Confirm   func(path string) (bool, error)
	Version   string
	CreatedAt time.Time
}

// OutputFile descreve um arquivo gerado.
type OutputFile struct {
	Path   string
	Format generator.Format
	Size   int64
}

// BuildResult é o resultado de uma geração.
type BuildResult struct {
	Files    []OutputFile
	Warnings []resume.Issue
	Lang     string
}

// Build valida o YAML e gera os formatos pedidos. Todas as saídas são
// verificadas antes de qualquer escrita, e cada arquivo é gravado de forma
// atômica (spec 002, FR-5, FR-6).
func Build(opts BuildOptions) (BuildResult, error) {
	var res BuildResult
	if len(opts.Formats) == 0 {
		opts.Formats = []generator.Format{generator.PDF}
	}

	r, rep, err := Load(opts.Input)
	res.Warnings = rep.Warnings
	if err != nil {
		return res, err
	}

	lang := opts.Lang
	if lang == "" {
		lang = r.Meta.Locale
	}
	labels, err := i18n.Get(lang)
	if err != nil {
		return res, err
	}
	res.Lang = labels.Code

	paths, err := OutputPaths(opts.Input, opts.Output, opts.Formats)
	if err != nil {
		return res, err
	}
	if err := checkOverwrite(paths, opts.Force, opts.Confirm); err != nil {
		return res, err
	}

	gopts := generator.Options{Labels: labels, Version: opts.Version, CreatedAt: opts.CreatedAt}
	for i, f := range opts.Formats {
		gen, err := generator.New(f, gopts)
		if err != nil {
			return res, err
		}
		var buf bytes.Buffer
		if err := gen.Generate(&buf, r); err != nil {
			return res, fmt.Errorf("erro ao gerar %s: %w", f, err)
		}
		if err := writeAtomic(paths[i], buf.Bytes()); err != nil {
			return res, err
		}
		res.Files = append(res.Files, OutputFile{Path: paths[i], Format: f, Size: int64(buf.Len())})
	}
	return res, nil
}

// OutputPaths calcula o caminho de cada formato (spec 002, FR-4, FR-7).
// Com um formato, output é o arquivo de saída (ou um diretório existente).
// Com vários, output é sempre um diretório.
func OutputPaths(input, output string, formats []generator.Format) ([]string, error) {
	base := strings.TrimSuffix(filepath.Base(input), filepath.Ext(input))
	dir := filepath.Dir(input)

	if len(formats) == 1 && output != "" {
		if isDirPath(output) {
			return []string{filepath.Join(output, base+formats[0].Ext())}, nil
		}
		return []string{output}, nil
	}
	if output != "" {
		if fi, err := os.Stat(output); err == nil && !fi.IsDir() {
			return nil, fmt.Errorf("com vários formatos, --output deve ser um diretório: %s", output)
		}
		dir = output
	}
	paths := make([]string, len(formats))
	for i, f := range formats {
		paths[i] = filepath.Join(dir, base+f.Ext())
	}
	return paths, nil
}

// isDirPath indica se output é um diretório: um que já existe ou um caminho
// terminado em separador ("dist/"), que será criado.
func isDirPath(output string) bool {
	if strings.HasSuffix(output, "/") || strings.HasSuffix(output, string(filepath.Separator)) {
		return true
	}
	fi, err := os.Stat(output)
	return err == nil && fi.IsDir()
}

func checkOverwrite(paths []string, force bool, confirm func(string) (bool, error)) error {
	if force {
		return nil
	}
	for _, p := range paths {
		if _, err := os.Stat(p); err != nil {
			continue
		}
		if confirm == nil {
			return fmt.Errorf("%w: %s (use --force para sobrescrever)", ErrOutputExists, p)
		}
		ok, err := confirm(p)
		if err != nil {
			return err
		}
		if !ok {
			return ErrCanceled
		}
	}
	return nil
}

// writeAtomic grava em um arquivo temporário no mesmo diretório e renomeia,
// para nunca deixar uma saída pela metade.
func writeAtomic(path string, data []byte) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	tmp, err := os.CreateTemp(dir, ".cv-craft-*")
	if err != nil {
		return err
	}
	defer func() { _ = os.Remove(tmp.Name()) }() // no-op após o rename
	if _, err := tmp.Write(data); err != nil {
		_ = tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	if err := os.Chmod(tmp.Name(), 0o644); err != nil {
		return err
	}
	return os.Rename(tmp.Name(), path)
}

// Init cria um YAML modelo. Não sobrescreve um arquivo existente sem force
// (spec 002, FR-9).
func Init(path string, force bool) error {
	if !force {
		if _, err := os.Stat(path); err == nil {
			return fmt.Errorf("%w: %s (use --force para sobrescrever)", ErrOutputExists, path)
		}
	}
	return writeAtomic(path, examples.Minimal)
}
