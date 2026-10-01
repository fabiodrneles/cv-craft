# Mesmos passos do CI (.github/workflows/ci.yml), para rodar localmente.

VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
COMMIT  ?= $(shell git rev-parse --short HEAD 2>/dev/null)
DATE    ?= $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS := -s -w -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)

.PHONY: all build test race cover lint smoke golden ci release-snapshot clean

all: ci

build: ## Compila o binário em ./bin
	go build -trimpath -ldflags "$(LDFLAGS)" -o bin/cv-craft .

test: ## Testes unitários, golden files e CLI
	go test -count=1 ./...

race: ## Testes com race detector
	go test -race -count=1 ./...

cover: ## Cobertura (mínimo de 80% em internal/)
	go test -count=1 -covermode=atomic -coverprofile=coverage.out ./internal/...
	scripts/coverage.sh coverage.out 80

lint: ## go vet + golangci-lint
	go vet ./...
	golangci-lint run ./...

smoke: ## Smoke test de ponta a ponta com o binário real
	scripts/smoke.sh

golden: ## Regrava os golden files após uma mudança intencional nas saídas
	go test ./internal/generator -update

ci: lint race cover smoke ## Tudo que o CI roda

release-snapshot: ## Gera os binários e arquivos da release em ./dist, sem publicar
	goreleaser release --snapshot --clean
	scripts/check-release.sh dist

clean:
	rm -rf bin dist coverage.out
