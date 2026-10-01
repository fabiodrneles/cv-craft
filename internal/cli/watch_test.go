package cli

import (
	"bytes"
	"context"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/fabiodrneles/cv-craft/examples"
)

// syncBuffer é um bytes.Buffer seguro para ler enquanto o --watch escreve.
type syncBuffer struct {
	mu  sync.Mutex
	buf bytes.Buffer
}

func (b *syncBuffer) Write(p []byte) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buf.Write(p)
}

func (b *syncBuffer) String() string {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buf.String()
}

type watchRun struct {
	out, err *syncBuffer
	cancel   context.CancelFunc
	done     chan int      // recebe o exit code
	exited   chan struct{} // fechado quando Run retorna
}

// startWatch roda "build <args> --watch" em segundo plano, com intervalos
// curtos para o teste.
func startWatch(t *testing.T, args ...string) *watchRun {
	t.Helper()
	ctx, cancel := context.WithCancel(context.Background())
	w := &watchRun{out: &syncBuffer{}, err: &syncBuffer{}, cancel: cancel, done: make(chan int, 1), exited: make(chan struct{})}
	go func() {
		defer close(w.exited)
		w.done <- Run(append(append([]string{"build"}, args...), "--watch"), Env{
			Stdout:       w.out,
			Stderr:       w.err,
			Context:      ctx,
			PollInterval: 10 * time.Millisecond,
			Debounce:     20 * time.Millisecond,
		})
	}()
	t.Cleanup(func() {
		cancel()
		select {
		case <-w.exited:
		case <-time.After(5 * time.Second):
			t.Error("o --watch não encerrou depois do cancelamento")
		}
	})
	return w
}

// waitFor espera até que cond seja verdadeira ou estoure o prazo.
func (w *watchRun) waitFor(t *testing.T, what string, limit time.Duration, cond func() bool) {
	t.Helper()
	deadline := time.Now().Add(limit)
	for !cond() {
		if time.Now().After(deadline) {
			t.Fatalf("%s não aconteceu em %v\nstdout: %s\nstderr: %s", what, limit, w.out, w.err)
		}
		time.Sleep(5 * time.Millisecond)
	}
}

// save grava o cv.yaml com uma data de modificação nova, para que a mudança
// seja percebida mesmo em sistemas de arquivos com resolução de 1 s ou 2 s.
func save(t *testing.T, content string, n int) {
	t.Helper()
	const path = "cv.yaml"
	mustWrite(t, path, content)
	mod := time.Now().Add(time.Duration(n) * 10 * time.Second)
	if err := os.Chtimes(path, mod, mod); err != nil {
		t.Fatal(err)
	}
}

// 002/AC-11 e AC-12: gera ao iniciar, regenera em menos de 1 s depois de
// salvar, mostra os erros de um YAML inválido sem encerrar e volta a gerar
// quando ele é corrigido; Ctrl+C (cancelamento) encerra com 0.
func TestWatchRebuildsOnSave(t *testing.T) {
	chdir(t)
	w := startWatch(t, "cv.yaml", "-f", "txt")
	generated := func(n int) func() bool {
		return func() bool { return strings.Count(w.out.String(), "Gerado: cv.txt") >= n }
	}
	w.waitFor(t, "primeira geração", 5*time.Second, generated(1))

	save(t, strings.Replace(string(examples.Minimal), `name: "Seu Nome"`, `name: "Ana Lima"`, 1), 1)
	w.waitFor(t, "regeneração depois de salvar", time.Second, generated(2))
	if txt, _ := os.ReadFile("cv.txt"); !strings.Contains(string(txt), "Ana Lima") {
		t.Fatalf("cv.txt não foi regenerado com o conteúdo novo:\n%s", txt)
	}

	save(t, "contact:\n  name: Ana\n", 2)
	w.waitFor(t, "erros do YAML inválido", time.Second, func() bool {
		return strings.Contains(w.err.String(), "problema(s)")
	})
	select {
	case code := <-w.done:
		t.Fatalf("o --watch encerrou com %d diante de um YAML inválido", code)
	default:
	}

	save(t, string(examples.Minimal), 3)
	w.waitFor(t, "geração depois da correção", time.Second, generated(3))

	w.cancel()
	select {
	case code := <-w.done:
		if code != ExitOK {
			t.Errorf("exit %d ao encerrar, esperava 0", code)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("o --watch não encerrou depois do cancelamento")
	}
}

// Arquivo de entrada inexistente: erro imediato (exit 1), sem ficar
// observando um caminho digitado errado.
func TestWatchMissingInput(t *testing.T) {
	chdir(t)
	w := startWatch(t, "nao-existe.yaml")
	select {
	case code := <-w.done:
		if code != ExitError {
			t.Errorf("exit %d, esperava %d", code, ExitError)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("o --watch deveria encerrar quando o arquivo não existe")
	}
}

// Saída existente sem --force fora de um terminal: nenhuma geração teria
// como gravar, então o --watch encerra com 3, como o build comum (002/FR-5).
func TestWatchOutputExists(t *testing.T) {
	chdir(t)
	mustWrite(t, "cv.pdf", "antigo")
	w := startWatch(t, "cv.yaml")
	select {
	case code := <-w.done:
		if code != ExitExists {
			t.Errorf("exit %d, esperava %d\nstderr: %s", code, ExitExists, w.err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("o --watch deveria encerrar quando a saída existe sem --force")
	}
	if b, _ := os.ReadFile("cv.pdf"); string(b) != "antigo" {
		t.Error("cv.pdf foi alterado")
	}
}

// Depois da primeira geração, as seguintes sobrescrevem as próprias saídas
// sem perguntar, mesmo sem --force.
func TestWatchOverwritesOwnOutput(t *testing.T) {
	chdir(t)
	w := startWatch(t, "cv.yaml", "-f", "md")
	w.waitFor(t, "primeira geração", 5*time.Second, func() bool {
		return strings.Contains(w.out.String(), "Gerado: cv.md")
	})
	save(t, string(examples.Minimal)+"\n", 1)
	w.waitFor(t, "segunda geração", time.Second, func() bool {
		return strings.Count(w.out.String(), "Gerado: cv.md") == 2
	})
	if strings.Contains(w.err.String(), "já existe") {
		t.Errorf("a segunda geração não deveria reclamar da própria saída:\n%s", w.err)
	}
}
