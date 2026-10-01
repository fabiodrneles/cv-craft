//go:build !windows

// Os testes deste arquivo rodam o binário real num pseudo-terminal (PTY),
// como uma pessoa no terminal: verificam a detecção de terminal, o Ctrl+D e o
// Ctrl+C de verdade, que os testes de internal/cli simulam (spec 005, AC-2 e
// AC-3; #15). O Windows não tem PTY no estilo Unix (o ConPTY é outra API), por
// isso o arquivo não é compilado lá; o modo interativo continua coberto pelos
// testes simulados em todos os sistemas.

package main

import (
	"bytes"
	"errors"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/creack/pty"
)

const ptyPrompt = "cv-craft›"

// ptySession é o binário rodando num PTY, com a saída acumulada.
type ptySession struct {
	t    *testing.T
	cmd  *exec.Cmd
	tty  *os.File
	mu   sync.Mutex
	out  bytes.Buffer
	read chan struct{} // fechado quando o PTY não tem mais saída
}

// buildBinary compila o cv-craft uma vez por teste, num diretório temporário.
func buildBinary(t *testing.T) string {
	t.Helper()
	bin := filepath.Join(t.TempDir(), "cv-craft")
	if out, err := exec.Command("go", "build", "-o", bin, ".").CombinedOutput(); err != nil {
		t.Fatalf("go build: %v\n%s", err, out)
	}
	return bin
}

// startUI roda "cv-craft ui" num PTY, num diretório de trabalho com uma cópia
// de examples/full.yaml.
func startUI(t *testing.T) (*ptySession, string) {
	t.Helper()
	bin := buildBinary(t)
	dir := t.TempDir()
	full, err := os.ReadFile("examples/full.yaml")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(dir, "examples"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "examples", "full.yaml"), full, 0o644); err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command(bin, "ui")
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "NO_COLOR=1", "SOURCE_DATE_EPOCH=1735787045")
	tty, err := pty.Start(cmd)
	if err != nil {
		t.Fatalf("pty.Start: %v", err)
	}
	s := &ptySession{t: t, cmd: cmd, tty: tty, read: make(chan struct{})}
	go func() {
		defer close(s.read)
		buf := make([]byte, 4096)
		for {
			n, err := tty.Read(buf)
			s.mu.Lock()
			s.out.Write(buf[:n])
			s.mu.Unlock()
			if err != nil {
				return // EIO quando o processo termina
			}
		}
	}()
	t.Cleanup(func() {
		_ = cmd.Process.Kill()
		_ = tty.Close()
	})
	return s, dir
}

func (s *ptySession) output() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.out.String()
}

// expect espera que text apareça na saída depois da posição from e devolve a
// posição logo após ele.
func (s *ptySession) expect(text string, from int) int {
	s.t.Helper()
	deadline := time.Now().Add(15 * time.Second)
	for {
		out := s.output()
		if from <= len(out) {
			if i := strings.Index(out[from:], text); i >= 0 {
				return from + i + len(text)
			}
		}
		if time.Now().After(deadline) {
			s.t.Fatalf("não apareceu %q na saída do PTY:\n%s", text, out)
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func (s *ptySession) send(text string) {
	s.t.Helper()
	if _, err := io.WriteString(s.tty, text); err != nil {
		s.t.Fatalf("escrever no PTY: %v", err)
	}
}

// wait espera o processo terminar e devolve o estado de saída.
func (s *ptySession) wait() *os.ProcessState {
	s.t.Helper()
	done := make(chan error, 1)
	go func() { done <- s.cmd.Wait() }()
	select {
	case err := <-done:
		var exit *exec.ExitError
		if err != nil && !errors.As(err, &exit) {
			s.t.Fatalf("wait: %v", err)
		}
		return s.cmd.ProcessState
	case <-time.After(10 * time.Second):
		s.t.Fatalf("o processo não terminou; saída:\n%s", s.output())
		return nil
	}
}

// 005/AC-2: num terminal real, build + exit gera o PDF e sai com 0.
func TestPTYBuildAndExit(t *testing.T) {
	s, dir := startUI(t)
	pos := s.expect(ptyPrompt, 0)
	s.send("build examples/full.yaml -o x.pdf\n")
	pos = s.expect("Gerado: x.pdf", pos)
	s.expect(ptyPrompt, pos)
	s.send("exit\n")
	if st := s.wait(); st.ExitCode() != 0 {
		t.Fatalf("exit %d, esperava 0; saída:\n%s", st.ExitCode(), s.output())
	}
	if fi, err := os.Stat(filepath.Join(dir, "x.pdf")); err != nil || fi.Size() == 0 {
		t.Fatalf("x.pdf não foi gerado: %v", err)
	}
}

// 005/AC-3: Ctrl+D no prompt encerra com 0.
func TestPTYCtrlD(t *testing.T) {
	s, _ := startUI(t)
	s.expect(ptyPrompt, 0)
	s.send("\x04") // Ctrl+D: fim da entrada no modo canônico do terminal
	if st := s.wait(); st.ExitCode() != 0 {
		t.Fatalf("exit %d, esperava 0; saída:\n%s", st.ExitCode(), s.output())
	}
}

// 005/FR-1: Ctrl+C encerra o modo interativo pelo sinal SIGINT, que o shell
// mostra como exit 130.
func TestPTYCtrlC(t *testing.T) {
	s, _ := startUI(t)
	s.expect(ptyPrompt, 0)
	s.send("\x03") // Ctrl+C: o terminal envia SIGINT ao processo
	st := s.wait()
	ws, ok := st.Sys().(syscall.WaitStatus)
	if !ok || !ws.Signaled() || ws.Signal() != syscall.SIGINT {
		t.Fatalf("esperava término por SIGINT (exit 130 no shell), obtive %v; saída:\n%s", st, s.output())
	}
}

// 005/FR-4 e 002/FR-5: num terminal, sobrescrever pergunta e respeita a
// resposta.
func TestPTYOverwritePrompt(t *testing.T) {
	s, dir := startUI(t)
	out := filepath.Join(dir, "x.pdf")
	if err := os.WriteFile(out, []byte("antigo"), 0o644); err != nil {
		t.Fatal(err)
	}

	pos := s.expect(ptyPrompt, 0)
	s.send("build examples/full.yaml -o x.pdf\n")
	pos = s.expect("Sobrescrever? [s/N]", pos)
	s.send("n\n")
	pos = s.expect("Cancelado", pos)
	if b, _ := os.ReadFile(out); string(b) != "antigo" {
		t.Fatal("x.pdf foi alterado apesar da resposta \"n\"")
	}

	pos = s.expect(ptyPrompt, pos)
	s.send("build examples/full.yaml -o x.pdf\n")
	pos = s.expect("Sobrescrever? [s/N]", pos)
	s.send("s\n")
	pos = s.expect("Gerado: x.pdf", pos)
	if b, _ := os.ReadFile(out); !bytes.HasPrefix(b, []byte("%PDF")) {
		t.Fatal("x.pdf não foi regravado depois da resposta \"s\"")
	}

	s.expect(ptyPrompt, pos)
	s.send("exit\n")
	if st := s.wait(); st.ExitCode() != 0 {
		t.Fatalf("exit %d, esperava 0", st.ExitCode())
	}
}
