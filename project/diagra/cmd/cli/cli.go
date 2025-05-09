package cli

import (
	"diagra/interpreter"
	"diagra/renderer"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	exampleDir = "example"
	outputDir  = "output"
)

// Run är ingångspunkten för CLI-programmet.
func RunCLI(args []string) {
	fs := flag.NewFlagSet("diagram", flag.ExitOnError)
	renderFile := fs.String("render", "", "Render a .diag file to SVG")

	fs.Parse(args)

	if *renderFile != "" {
		fmt.Println("Rendering:", *renderFile)
		// call your rendering logic
		return
	}

	fs.Usage()
}

func renderAllDiagrams(diagramFiles []string) string {

	files, err := os.ReadDir(exampleDir)
	if err != nil {
		fmt.Println("Kunde inte läsa katalog:", err)
	}

	for _, file := range files {
		if !strings.HasSuffix(file.Name(), ".diag") {
			continue
		}
		fullPath := filepath.Join(exampleDir, file.Name())
		processDiagramFile(fullPath, outputDir)
	}
	return "Rendering finished"
}

func processDiagramFile(path string, outputDir string) {
	fmt.Println("Läser:", path)

	src, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("Kunde inte läsa fil:", err)
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
		fmt.Println("Tolkningsfel:", err)
		return
	}

	svg := renderer.RenderSVG(diagram)

	base := strings.TrimSuffix(filepath.Base(path), ".diag")
	outPath := filepath.Join(outputDir, base+".svg")

	err = os.WriteFile(outPath, []byte(svg), 0644)
	if err != nil {
		fmt.Println("Kunde inte spara SVG:", err)
		return
	}

	fmt.Println("Skapade:", outPath)
}
