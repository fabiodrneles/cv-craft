package cli

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"cv-craft/internal/generator"
	"cv-craft/internal/parser"

	"github.com/fatih/color"
)

type GeminiUI struct {
	scanner *bufio.Scanner
}

func NewGeminiUI() *GeminiUI {
	return &GeminiUI{
		scanner: bufio.NewScanner(os.Stdin),
	}
}

// ShowHeader moderno, clean e minimalista
func (ui *GeminiUI) ShowHeader() {
	cyan := color.New(color.FgCyan)
	blue := color.New(color.FgBlue)
	boldCyan := cyan.Add(color.Bold)
	boldBlue := blue.Add(color.Bold)

	fmt.Println(boldCyan.Sprint("┌────────────────────────────────────────────────────────┐"))
	fmt.Println(boldCyan.Sprint("│                                                        │"))
	fmt.Println(boldCyan.Sprint("│  ") + boldBlue.Sprint("● CV-CRAFT CLI - Professional Resume Generator") + boldCyan.Sprint("        │"))
	fmt.Println(boldCyan.Sprint("│                                                        │"))
	fmt.Println(boldCyan.Sprint("└────────────────────────────────────────────────────────┘"))
	fmt.Println()

	ui.ShowInfo(" Professional Resume Generator built with Go")
	ui.ShowInfo(" Multi-format: PDF, Markdown, and Text")
	ui.ShowInfo(" ATS-optimized for better job applications")
	ui.ShowInfo(" Type 'help' for commands or 'exit' to quit")
	fmt.Println()
}

// Loading animation similar to Gemini CLI
func (ui *GeminiUI) ShowLoading(message string) {
	frames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	blue := color.New(color.FgBlue)

	fmt.Printf("\r%s %s", blue.Sprint("│"), message)

	for i := 0; i < 10; i++ {
		fmt.Printf("\r%s %s %s", blue.Sprint("│"), blue.Sprint(frames[i%len(frames)]), message)
		time.Sleep(100 * time.Millisecond)
	}
	fmt.Printf("\r%s ✓ %s\n", blue.Sprint("│"), message)
}

// Interactive prompt com lambda (λ) - o "y invertido"
func (ui *GeminiUI) Prompt(question string) string {
	color.New(color.FgCyan)
	blue := color.New(color.FgBlue)
	yellow := color.New(color.FgYellow)

	fmt.Printf("\n%s %s\n", blue.Sprint("│"), question)
	fmt.Printf("%s ", yellow.Sprint("λ")) // Lambda symbol instead of $

	ui.scanner.Scan()
	return strings.TrimSpace(ui.scanner.Text())
}

// Show info in Gemini style
func (ui *GeminiUI) ShowInfo(message string) {
	blue := color.New(color.FgBlue)
	fmt.Printf("%s ℹ %s\n", blue.Sprint("│"), message)
}

// Show success message
func (ui *GeminiUI) ShowSuccess(message string) {
	green := color.New(color.FgGreen)
	fmt.Printf("%s ✓ %s\n", green.Sprint("│"), message)
}

// Show error message
func (ui *GeminiUI) ShowError(message string) {
	red := color.New(color.FgRed)
	fmt.Printf("%s ✗ %s\n", red.Sprint("│"), message)
}

// Show warning message
func (ui *GeminiUI) ShowWarning(message string) {
	yellow := color.New(color.FgYellow)
	fmt.Printf("%s ⚠ %s\n", yellow.Sprint("│"), message)
}

// Show command output
func (ui *GeminiUI) ShowCommandOutput(output string) {
	blue := color.New(color.FgBlue)
	fmt.Printf("%s %s\n", blue.Sprint("│"), output)
}

// Show section divider
func (ui *GeminiUI) ShowDivider() {
	cyan := color.New(color.FgCyan)
	fmt.Println(cyan.Sprint("├────────────────────────────────────────────────────────┤"))
}

// Show help command com layout melhorado
func (ui *GeminiUI) ShowHelp() {
	cyan := color.New(color.FgCyan)
	green := color.New(color.FgGreen)
	yellow := color.New(color.FgYellow)
	blue := color.New(color.FgBlue)

	ui.ShowDivider()
	fmt.Println()
	fmt.Println(cyan.Sprint("📖 Available Commands:"))
	fmt.Println()

	commands := []struct {
		cmd     string
		desc    string
		example string
	}{
		{"build", "Generate resume from YAML", "build curriculum.yaml --format pdf"},
		{"init", "Create a new resume template", "init my-resume.yaml"},
		{"validate", "Validate YAML structure", "validate curriculum.yaml"},
		{"templates", "List available templates", "templates list"},
		{"help", "Show this help message", "help"},
		{"exit", "Exit the CLI", "exit"},
	}

	for _, cmd := range commands {
		fmt.Printf("  %s %s\n", green.Sprint(cmd.cmd), cmd.desc)
		if cmd.example != "" {
			fmt.Printf("     %sExample: %s\n", blue.Sprint("│ "), yellow.Sprint(cmd.example))
		}
		fmt.Println()
	}

	ui.ShowDivider()
	ui.ShowInfo("Flags: --format (pdf|md|text), --output, --ats")
	ui.ShowInfo("Tips: Use TAB for auto-completion, ↑↓ for command history")
	ui.ShowDivider()
}

// Interactive mode
func (ui *GeminiUI) InteractiveMode() {
	ui.ShowHeader()

	for {
		input := ui.Prompt("cv-craft")

		if input == "exit" || input == "quit" {
			ui.ShowSuccess("Thanks for using CV-Craft!")
			break
		}

		if input == "help" || input == "" {
			ui.ShowHelp()
			continue
		}

		if input == "clear" || input == "cls" {
			ui.clearScreen()
			continue
		}

		if input == "version" || input == "v" {
			ui.ShowInfo("CV-Craft CLI v1.0 • Built with Go • ATS Optimized")
			continue
		}

		// Parse and execute commands
		ui.executeCommand(input)
	}
}

func (ui *GeminiUI) clearScreen() {
	fmt.Print("\033[H\033[2J")
	ui.ShowHeader()
}

func (ui *GeminiUI) executeCommand(input string) {
	args := strings.Fields(input)
	if len(args) == 0 {
		return
	}

	command := args[0]

	switch command {
	case "build":
		ui.handleBuildCommand(args[1:])
	case "init":
		ui.handleInitCommand(args[1:])
	case "validate":
		ui.handleValidateCommand(args[1:])
	case "templates":
		ui.handleTemplatesCommand(args[1:])
	default:
		ui.ShowError(fmt.Sprintf("Unknown command: %s. Type 'help' for available commands.", command))
	}
}

// Função corrigida para geração real de PDF no modo interativo
func (ui *GeminiUI) handleBuildCommand(args []string) {
	if len(args) == 0 {
		ui.ShowError("Please specify a YAML file. Example: build curriculum.yaml")
		return
	}

	inputFile := args[0]
	format := "pdf"
	outputFile := ""
	ats := false

	// Parse flags corretamente
	for i := 1; i < len(args); i++ {
		switch args[i] {
		case "--format", "-format":
			if i+1 < len(args) {
				format = strings.ToLower(args[i+1])
				i++
			}
		case "--output", "-o", "-output":
			if i+1 < len(args) {
				outputFile = args[i+1]
				i++
			}
		case "--ats", "-ats":
			ats = true
		}
	}

	ui.ShowLoading(fmt.Sprintf("Generating resume in %s format...", strings.ToUpper(format)))

	// Chamar a lógica REAL de geração
	if err := ui.executeBuild(inputFile, outputFile, format, ats); err != nil {
		ui.ShowError(fmt.Sprintf("Failed to generate resume: %v", err))
		return
	}

	// Definir arquivo de saída para exibição
	if outputFile == "" {
		baseName := strings.TrimSuffix(inputFile, filepath.Ext(inputFile))
		ext := format
		switch format {
		case "md", "markdown":
			ext = "md"
		case "text":
			ext = "txt"
		}
		outputFile = fmt.Sprintf("%s.%s", baseName, ext)
	}

	// Verificar se o arquivo foi realmente criado
	if fileInfo, err := os.Stat(outputFile); err == nil {
		ui.ShowSuccess(fmt.Sprintf("Resume generated successfully: %s (%d bytes)", outputFile, fileInfo.Size()))

		// Mensagem adicional baseada no formato
		switch format {
		case "pdf":
			ui.ShowInfo("🚀 Ready to share with recruiters! 📄")
		case "md", "markdown":
			ui.ShowInfo("📝 Markdown file generated - perfect for online portfolios")
		case "text":
			ui.ShowInfo("📄 Text file generated - compatible with any system")
		}
	} else {
		ui.ShowWarning("Generation completed but could not verify output file")
		ui.ShowInfo(fmt.Sprintf("Expected file: %s", outputFile))
	}
}

// Função que executa a geração REAL (similar à buildResume do main.go)
func (ui *GeminiUI) executeBuild(inputFile, outputFile, format string, ats bool) error {
	// Verificar se o arquivo de entrada existe
	if _, err := os.Stat(inputFile); os.IsNotExist(err) {
		return fmt.Errorf("input file not found: %s", inputFile)
	}

	// Validar formato
	validFormats := map[string]bool{
		"pdf":      true,
		"markdown": true,
		"md":       true,
		"text":     true,
	}
	if !validFormats[format] {
		return fmt.Errorf("invalid format: %s. Use: pdf, markdown, md or text", format)
	}

	// Parse do arquivo YAML
	cv, err := parser.ParseFile(inputFile)
	if err != nil {
		return fmt.Errorf("error processing YAML file: %v", err)
	}

	// Validação do currículo
	if err := parser.ValidateCurriculum(cv); err != nil {
		return fmt.Errorf("curriculum validation error: %v", err)
	}

	// Definir arquivo de saída padrão
	outputPath := outputFile
	if outputPath == "" {
		baseName := strings.TrimSuffix(inputFile, filepath.Ext(inputFile))
		ext := format
		switch format {
		case "md", "markdown":
			ext = "md"
		case "text":
			ext = "txt"
		}
		outputPath = fmt.Sprintf("%s.%s", baseName, ext)
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
		return fmt.Errorf("generator not implemented for format: %s", format)
	}

	// Gerar o currículo
	if err := gen.Generate(cv, outputPath); err != nil {
		return fmt.Errorf("error generating resume: %v", err)
	}

	return nil
}

func (ui *GeminiUI) handleInitCommand(args []string) {
	if len(args) == 0 {
		ui.ShowError("Please specify a filename. Example: init my-resume.yaml")
		return
	}

	filename := args[0]
	ui.ShowLoading("Creating new resume template...")

	// Create a basic template
	template := `contact:
  name: "Your Name"
  email: "your.email@example.com"
  phone: "+55 (11) 99999-9999"
  linkedin: "linkedin.com/in/yourprofile"
  portfolio: "yourportfolio.com"
  location: "City, Country"

professional_title: "Your Professional Title"

summary: "Brief professional summary highlighting your experience and skills."

skills:
  - name: "Category 1"
    level: "proficient"
    keywords: ["Skill1", "Skill2", "Skill3"]
  
  - name: "Category 2" 
    level: "intermediate"
    keywords: ["Skill4", "Skill5"]

experience:
  - role: "Your Role"
    company: "Company Name"
    location: "City, Country"
    period: "2023 - Present"
    responsibilities:
      - "Describe your responsibilities and achievements"
      - "Use metrics and results when possible"
    technologies: ["Tech1", "Tech2", "Tech3"]

education:
  - degree: "Your Degree"
    institution: "University Name"
    location: "City, Country"
    period: "2019 - 2023"
`

	err := os.WriteFile(filename, []byte(template), 0644)
	if err != nil {
		ui.ShowError(fmt.Sprintf("Error creating file: %v", err))
		return
	}

	ui.ShowSuccess(fmt.Sprintf("Template created: %s", filename))
	ui.ShowInfo("Edit the file with your information and run: build " + filename)
}

func (ui *GeminiUI) handleValidateCommand(args []string) {
	if len(args) == 0 {
		ui.ShowError("Please specify a YAML file to validate")
		return
	}

	inputFile := args[0]
	ui.ShowLoading(fmt.Sprintf("Validating %s...", inputFile))

	// Validação REAL do arquivo YAML
	cv, err := parser.ParseFile(inputFile)
	if err != nil {
		ui.ShowError(fmt.Sprintf("Validation failed: %v", err))
		return
	}

	// Validação da estrutura
	if err := parser.ValidateCurriculum(cv); err != nil {
		ui.ShowError(fmt.Sprintf("Validation failed: %v", err))
		return
	}

	ui.ShowSuccess("YAML structure is valid! ✅")
	ui.ShowInfo(fmt.Sprintf("✓ Contact: %s", cv.Contact.Name))
	ui.ShowInfo(fmt.Sprintf("✓ Skills: %d categories", len(cv.Skills)))
	ui.ShowInfo(fmt.Sprintf("✓ Experience: %d positions", len(cv.Experience)))
	ui.ShowInfo(fmt.Sprintf("✓ Education: %d entries", len(cv.Education)))
}

func (ui *GeminiUI) handleTemplatesCommand(args []string) {
	if len(args) > 0 && args[0] == "list" {
		ui.ShowLoading("Loading available templates...")

		templates := []string{
			"📄 ATS Professional (Default)",
			"🎨 Modern Creative",
			"📊 Executive Summary",
			"🔬 Academic Research",
			"💻 Tech Specialist",
		}

		ui.ShowCommandOutput("Available Templates:")
		for _, template := range templates {
			ui.ShowCommandOutput("  • " + template)
		}
		ui.ShowInfo("Use: init filename.yaml to create a basic template")
	} else {
		ui.ShowError("Use: templates list")
	}
}
