package cli

import (
	"diagra/interpreter"
	"diagra/renderer"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	exampleDir = "example"
	outputDir  = "output"
)

var renderStart time.Time
var combinedTime int64

// Run is the entry point for the CLI application.
func RunCLI(args []string) {
	switch args[0] {
	case "render":
		if len(args) < 2 {
			fmt.Println("Specify a .diagra file to render")
			return
		}
		if !strings.HasSuffix(args[1], ".diag") {
			fmt.Println("File must have .diag extension")
			return
		}
		fullPath := filepath.Join(exampleDir, args[1])
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			fmt.Println("File does not exist:", fullPath)
			return
		}
		renderStart = time.Now()
		processDiagramFile(fullPath, outputDir)
		return
	case "render-all":
		renderStart = time.Now()
		renderAllDiagrams()
		return
	case "-h", "--help", "help":
		help()
		return
	default:
		fmt.Println("Run -h, --help or help for usage")
		return
	}

}

// renderAllDiagrams läser alla .diag-filer i exempel-katalogen
// och renderar dem till SVG-filer i output-katalogen.
func renderAllDiagrams() {

	files, err := os.ReadDir(exampleDir)
	if err != nil {
		fmt.Println("Could not read directory:", err)
	}

	for _, file := range files {
		if !strings.HasSuffix(file.Name(), ".diag") {
			continue
		}
		fullPath := filepath.Join(exampleDir, file.Name())
		processDiagramFile(fullPath, outputDir)
		combinedTime += time.Since(renderStart).Milliseconds()
	}
	fmt.Println("All diagrams rendered to SVG in", outputDir)
	fmt.Printf("Total time: %d ms\n", combinedTime)
}

// processDiagramFile läser en .diag-fil, tolkar den och
// sparar den som en SVG-fil i output-katalogen.
func processDiagramFile(path string, outputDir string) {
	fmt.Println("Reading:", path)

	src, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("Could not read file:", err)
		return
	}
	fmt.Printf("Input length: %d bytes\n", len(src))
	tokens := interpreter.Lex(string(src))
	for i, tok := range tokens {
		fmt.Printf("%d: %s (%s)\n", i, tok.Value, tok.Type)
	}
	diagram, err := interpreter.Parse(tokens)
	fmt.Printf("Nodes: %d, Edges: %d\n", len(diagram.Nodes), len(diagram.Edges))
	if err != nil {
		fmt.Println("Parsing error:", err)
		return
	}

	svg := renderer.RenderSVG(diagram)

	base := strings.TrimSuffix(filepath.Base(path), ".diag")
	outPath := filepath.Join(outputDir, base+".svg")

	err = os.WriteFile(outPath, []byte(svg), 0644)
	if err != nil {
		fmt.Println("Could not save SVG:", err)
		return
	}
	renderTime := time.Since(renderStart).Milliseconds()
	fmt.Println("Created:", outPath)
	fmt.Printf("Render time: %d ms\n", renderTime)
}

func help() {
	fmt.Println("Usage: diagra [command]")
	fmt.Println("Commands:")
	fmt.Println("  render <file>		Render a diagram from a .diag file")
	fmt.Println("  render-all		Render all diagrams in the example directory")
	fmt.Println("  -h, --help, help     	Show this help message")
}
