// Comando cv-craft gera currículos profissionais e otimizados para ATS a
// partir de um arquivo YAML.
package main

import (
	"os"
	"runtime/debug"

	"github.com/mattn/go-isatty"

	"github.com/fabiodrneles/cv-craft/internal/cli"
)

// Preenchidos via -ldflags no build de release.
var (
	version = "dev"
	commit  = ""
	date    = ""
)

func main() {
	os.Exit(cli.Run(os.Args[1:], cli.Env{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		IsTerminal: func() bool {
			fd := os.Stdin.Fd()
			return isatty.IsTerminal(fd) || isatty.IsCygwinTerminal(fd)
		},
		Getenv:  os.Getenv,
		Version: resolveVersion(),
		Commit:  commit,
		Date:    date,
	}))
}

// resolveVersion usa a versão do módulo quando instalado com
// "go install ...@vX.Y.Z" e não há -ldflags.
func resolveVersion() string {
	if version != "dev" {
		return version
	}
	if bi, ok := debug.ReadBuildInfo(); ok && bi.Main.Version != "" && bi.Main.Version != "(devel)" {
		return bi.Main.Version
	}
	return version
}
