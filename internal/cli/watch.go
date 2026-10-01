package cli

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fabiodrneles/cv-craft/internal/app"
)

// Intervalos padrão do --watch (spec 002, FR-15): a data de modificação do
// YAML é conferida a cada pollInterval, e a geração espera o arquivo ficar
// estável por debounce, para não gerar no meio de um salvamento.
const (
	defaultPollInterval = 200 * time.Millisecond
	defaultDebounce     = 300 * time.Millisecond
)

// fileStamp identifica uma versão do arquivo observado.
type fileStamp struct {
	mod  time.Time
	size int64
}

func stampOf(path string) (fileStamp, bool) {
	fi, err := os.Stat(path)
	if err != nil {
		return fileStamp{}, false
	}
	return fileStamp{fi.ModTime(), fi.Size()}, true
}

// watch gera uma vez e depois regenera a cada alteração do YAML, até
// Ctrl+C (exit 0). Erros de validação ou de geração são mostrados sem
// encerrar; só a recusa em sobrescrever uma saída existente encerra, porque
// nenhuma geração seguinte teria como gravar.
func (r *runner) watch(opts app.BuildOptions, quiet bool) int {
	if _, err := os.Stat(opts.Input); err != nil {
		return r.fail(err)
	}
	ctx, stop := signal.NotifyContext(r.env.Context, os.Interrupt, syscall.SIGTERM)
	defer stop()

	poll := r.env.PollInterval
	if poll <= 0 {
		poll = defaultPollInterval
	}
	debounce := r.env.Debounce
	if debounce <= 0 {
		debounce = defaultDebounce
	}

	if !quiet {
		fmt.Fprintf(r.err, "Observando %s. Salve o arquivo para gerar de novo; Ctrl+C para sair.\n", opts.Input)
	}
	last, _ := stampOf(opts.Input)
	if code, stopWatching := r.watchBuild(&opts, quiet); stopWatching {
		return code
	}

	ticker := time.NewTicker(poll)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			if !quiet {
				fmt.Fprintln(r.err, "Observação encerrada.")
			}
			return ExitOK
		case <-ticker.C:
		}
		cur, ok := stampOf(opts.Input)
		if !ok || cur == last {
			// Editores que salvam renomeando um arquivo temporário deixam o
			// YAML ausente por um instante; espera ele voltar.
			continue
		}
		// Espera o arquivo parar de mudar antes de gerar.
		select {
		case <-ctx.Done():
			continue
		case <-time.After(debounce):
		}
		if again, ok := stampOf(opts.Input); !ok || again != cur {
			continue
		}
		last = cur
		if !quiet {
			fmt.Fprintf(r.err, "\n%s alterado; gerando de novo…\n", opts.Input)
		}
		if code, stopWatching := r.watchBuild(&opts, quiet); stopWatching {
			return code
		}
	}
}

// watchBuild faz uma geração do --watch. Depois da primeira geração bem
// sucedida, as seguintes sobrescrevem as próprias saídas sem perguntar.
func (r *runner) watchBuild(opts *app.BuildOptions, quiet bool) (code int, stopWatching bool) {
	res, err := app.Build(*opts)
	if !quiet {
		r.warnings(res.Warnings)
	}
	if err != nil {
		code := r.fail(err)
		if errors.Is(err, app.ErrOutputExists) || errors.Is(err, app.ErrCanceled) {
			return code, true
		}
		if !quiet {
			fmt.Fprintln(r.err, "Corrija o arquivo e salve de novo; a geração é refeita automaticamente.")
		}
		return 0, false
	}
	opts.Force = true
	if !quiet {
		for _, f := range res.Files {
			fmt.Fprintf(r.out, "Gerado: %s (%s)\n", f.Path, humanSize(f.Size))
		}
	}
	return 0, false
}
