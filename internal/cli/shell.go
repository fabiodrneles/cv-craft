package cli

import (
	"errors"
	"fmt"
	"strings"

	"github.com/fatih/color"
)

// interactive executa o modo interativo (spec 005). Cada linha é tratada como
// argumentos da CLI e executada pelo mesmo dispatch.
func (r *runner) interactive() int {
	if !r.env.IsTerminal() {
		r.errorf("o modo interativo requer um terminal; use os comandos diretamente (cv-craft help)")
		return ExitUsage
	}

	accent := color.New(color.FgCyan, color.Bold)
	fmt.Fprintf(r.out, "%s %s — modo interativo\n", accent.Sprint("CV-Craft"), r.env.Version)
	fmt.Fprintln(r.out, "Digite 'help' para ver os comandos ou 'exit' para sair.")

	for {
		fmt.Fprint(r.out, "\n"+accent.Sprint("cv-craft›")+" ")
		line, err := r.in.ReadString('\n')
		if err != nil && line == "" { // EOF (Ctrl+D) ou erro de leitura
			fmt.Fprintln(r.out)
			return ExitOK
		}
		line = strings.TrimSpace(line)
		switch line {
		case "":
			continue
		case "exit", "quit", "sair":
			return ExitOK
		case "clear", "cls":
			fmt.Fprint(r.out, "\033[H\033[2J")
			continue
		}

		args, perr := tokenize(line)
		if perr != nil {
			r.errorf("%v", perr)
			continue
		}
		if args[0] == "ui" || args[0] == "interactive" {
			fmt.Fprintln(r.out, "Você já está no modo interativo.")
			continue
		}
		r.dispatch(args)
	}
}

// tokenize separa uma linha em argumentos, respeitando aspas simples e
// duplas. Barras invertidas não são escape, para aceitar caminhos do Windows.
func tokenize(line string) ([]string, error) {
	var (
		args    []string
		cur     strings.Builder
		quote   rune
		started bool
	)
	for _, c := range line {
		switch {
		case quote != 0:
			if c == quote {
				quote = 0
			} else {
				cur.WriteRune(c)
			}
		case c == '"' || c == '\'':
			quote, started = c, true
		case c == ' ' || c == '\t':
			if started {
				args = append(args, cur.String())
				cur.Reset()
				started = false
			}
		default:
			cur.WriteRune(c)
			started = true
		}
	}
	if quote != 0 {
		return nil, errors.New("aspas não fechadas")
	}
	if started {
		args = append(args, cur.String())
	}
	if len(args) == 0 {
		return nil, errors.New("comando vazio")
	}
	return args, nil
}
