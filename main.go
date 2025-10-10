package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"cv-craft/cli"
	"cv-craft/internal/generator"
	"cv-craft/internal/parser"
)

func main() {
	// If no arguments, start interactive mode
	if len(os.Args) == 1 {
		ui := cli.NewGeminiUI()
		ui.InteractiveMode()
		return
	}

	// If first argument is "interactive" or "ui", start interactive mode
	if os.Args[1] == "interactive" || os.Args[1] == "ui" {
		ui := cli.NewGeminiUI()
		ui.InteractiveMode()
		return
	}

	// If using go run main.go with flags (legacy mode)
	if len(os.Args) > 1 && strings.Contains(os.Args[0], "main.go") {
		runLegacyMode()
		return
	}

	// Normal build command mode
	if len(os.Args) < 3 {
		printUsage()
		os.Exit(1)
	}

	// Verificar se é o comando "build"
	if os.Args[1] != "build" {
		printUsage()
		os.Exit(1)
	}

	// Processar argumentos manualmente para melhor controle
	inputFile := ""
	outputFile := ""
	format := "pdf"
	ats := false

	args := os.Args[2:] // Pular "cv-craft" e "build"

	// Processar flags e argumentos
	for i := 0; i < len(args); i++ {
		arg := args[i]

		switch arg {
		case "-format", "--format":
			if i+1 < len(args) {
				format = strings.ToLower(args[i+1])
				i++ // Pular o próximo argumento
			} else {
				log.Fatal("Flag -format requer um valor")
			}
		case "-o", "--output", "-output":
			if i+1 < len(args) {
				outputFile = args[i+1]
				i++ // Pular o próximo argumento
			} else {
				log.Fatal("Flag -o requer um valor")
			}
		case "-ats", "--ats":
			ats = true
		case "-v", "--verbose":
			log.SetFlags(log.LstdFlags | log.Lshortfile)
		default:
			// Se não é uma flag, deve ser o arquivo de entrada
			if strings.HasPrefix(arg, "-") {
				log.Fatalf("Flag desconhecida: %s", arg)
			} else if inputFile == "" {
				inputFile = arg
			} else {
				log.Fatalf("Argumento inesperado: %s", arg)
			}
		}
	}

	// Verificar se o arquivo de entrada foi fornecido
	if inputFile == "" {
		log.Fatal("Arquivo de entrada não especificado. Use: cv-craft build <arquivo.yaml>")
	}

	// Executar a construção do currículo
	if err := buildResume(inputFile, outputFile, format, ats); err != nil {
		log.Fatalf("❌ Erro ao gerar currículo: %v", err)
	}
}

func runLegacyMode() {
	// Modo legado para go run main.go com flags simples
	inputFile := ""
	outputFile := ""
	format := "pdf"
	ats := false

	for i := 1; i < len(os.Args); i++ {
		arg := os.Args[i]

		switch arg {
		case "-input":
			if i+1 < len(os.Args) {
				inputFile = os.Args[i+1]
				i++
			}
		case "-format":
			if i+1 < len(os.Args) {
				format = strings.ToLower(os.Args[i+1])
				i++
			}
		case "-output":
			if i+1 < len(os.Args) {
				outputFile = os.Args[i+1]
				i++
			}
		case "-ats":
			ats = true
		case "-v", "--verbose":
			log.SetFlags(log.LstdFlags | log.Lshortfile)
		case "-help", "--help":
			printUsage()
			os.Exit(0)
		}
	}

	if inputFile == "" {
		log.Fatal("Arquivo de entrada não especificado. Use: go run main.go -input curriculum.yaml")
	}

	if err := buildResume(inputFile, outputFile, format, ats); err != nil {
		log.Fatalf("❌ Erro ao gerar currículo: %v", err)
	}
}

func buildResume(inputFile, outputFile, format string, ats bool) error {
	// Verificar se o arquivo de entrada existe
	if _, err := os.Stat(inputFile); os.IsNotExist(err) {
		return fmt.Errorf("arquivo de entrada não encontrado: %s", inputFile)
	}

	// Validar formato
	validFormats := map[string]bool{
		"pdf":      true,
		"markdown": true,
		"md":       true,
		"text":     true,
	}
	if !validFormats[format] {
		return fmt.Errorf("formato inválido: %s. Use: pdf, markdown, md ou text", format)
	}

	// Parse do arquivo YAML
	cv, err := parser.ParseFile(inputFile)
	if err != nil {
		return fmt.Errorf("erro ao processar arquivo YAML: %v", err)
	}

	// Validação do currículo
	if err := parser.ValidateCurriculum(cv); err != nil {
		return fmt.Errorf("erro de validação do currículo: %v", err)
	}

	// Definir arquivo de saída padrão se não especificado
	outputPath := outputFile
	if outputPath == "" {
		baseName := strings.TrimSuffix(inputFile, filepath.Ext(inputFile))
		// Usar a extensão correta baseada no formato
		ext := format
		switch format {
		case "md", "markdown":
			ext = "md"
		case "text":
			ext = "txt"
		}
		outputPath = fmt.Sprintf("%s.%s", baseName, ext)
	}

	// VERIFICAR SE O ARQUIVO JÁ EXISTE
	if _, err := os.Stat(outputPath); err == nil {
		// Arquivo já existe, perguntar se quer sobrescrever?
		fmt.Printf("⚠️  Arquivo já existe: %s\n", outputPath)
		fmt.Print("   Deseja sobrescrever? (s/N): ")

		var response string
		fmt.Scanln(&response)

		if strings.ToLower(response) != "s" && strings.ToLower(response) != "sim" {
			return fmt.Errorf("geração cancelada pelo usuário")
		}
	}

	// Criar gerador baseado no formato
	var gen generator.Generator

	switch format {
	case "pdf":
		gen = &generator.PDFGenerator{IsATS: ats}
	case "markdown", "md":
		gen = &generator.MarkdownGenerator{}
	case "text":
		gen = &generator.TextGenerator{}
	default:
		return fmt.Errorf("gerador não implementado para o formato: %s", format)
	}

	// Gerar o currículo
	if err := gen.Generate(cv, outputPath); err != nil {
		return fmt.Errorf("erro ao gerar currículo: %v", err)
	}

	// VERIFICAÇÃO FINAL: Confirmar que o arquivo foi criado com sucesso
	fileInfo, err := os.Stat(outputPath)
	if err != nil {
		return fmt.Errorf("erro: arquivo de saída não foi criado - %v", err)
	}

	if fileInfo.Size() == 0 {
		return fmt.Errorf("aviso: arquivo de saída está vazio")
	}

	fmt.Printf("✅ Currículo gerado com sucesso: %s (%d bytes)\n", outputPath, fileInfo.Size())
	fmt.Printf("📁 Localização: %s\n", getFullPath(outputPath))

	// Mensagem adicional baseada no formato
	switch format {
	case "pdf":
		fmt.Println("🚀 Pronto para compartilhar com recrutadores! 📄")
	case "md", "markdown":
		fmt.Println("📝 Arquivo Markdown gerado - ideal para portfólios online")
	case "text":
		fmt.Println("📄 Arquivo de texto gerado - compatível com qualquer sistema")
	}

	return nil
}

// Função auxiliar para obter o caminho completo
func getFullPath(filename string) string {
	absPath, err := filepath.Abs(filename)
	if err != nil {
		return filename
	}
	return absPath
}

func printUsage() {
	fmt.Println(`GO-Craft - Professional Resume Generator

Uso:
  cv-craft build <arquivo.yaml> [flags]
  cv-craft interactive          (Interactive mode)
  cv-craft ui                   (Interactive mode)

Flags:
  -format string    Output format: pdf, markdown, md, text (default "pdf")
  -o string         Output file path (optional)
  -ats              Optimize for ATS - PDF only
  -v, --verbose     Verbose output for debugging

Exemplos:
  cv-craft build curriculum.yaml
  cv-craft build curriculum.yaml -format md
  cv-craft build curriculum.yaml -format text
  cv-craft build curriculum.yaml -format pdf -ats
  cv-craft build curriculum.yaml -o my_resume.pdf -format pdf -v

Interactive commands:
  cv-craft interactive    # Start interactive interface
  cv-craft ui             # Start interactive interface`)
}
