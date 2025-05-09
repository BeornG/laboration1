package tui

import (
	"diagra/interpreter"
	"diagra/renderer"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// Ingångspunkten för TUI-programmet
func RunTUI() {
	diagFiles, err := loadDiagFiles("./example")
	if err != nil {
		fmt.Println("Error loading .diag files:", err)
		os.Exit(1)
	}

	p := tea.NewProgram(InitialModel(diagFiles), tea.WithAltScreen())
	if err := p.Start(); err != nil { // depricated?
		fmt.Println("Error running TUI:", err)
		os.Exit(1)
	}
}

// loadingDiagFiles ladar alla .diag-filer i den angivna katalogen
// och returnerar en lista med filnamn.
func loadDiagFiles(path string) ([]string, error) {
	var files []string
	err := filepath.WalkDir(path, func(p string, d fs.DirEntry, e error) error {
		if e != nil {
			return e
		}
		if !d.IsDir() && filepath.Ext(p) == ".diag" {
			files = append(files, filepath.Base(p))
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, nil
}

// renderDiagToSVG läser en .diag-fil, tolkar den och
// sparar den som en SVG-fil i output-katalogen.
// Den returnerar sökvägen till den skapade SVG-filen.
func renderDiagToSVG(path string) string {
	outputDir := "output"
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		err := os.Mkdir(outputDir, 0755)
		if err != nil {
			fmt.Println("Could not create output directory:", err)
			return ""
		}
	}

	// fmt.Print("Reading:", path)
	src, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("Could not read file:", err)
		return ""
	}
	// fmt.Printf("Input length: %d bytes\n", len(src))

	tokens := interpreter.Lex(string(src))
	// for i, tok := range tokens {
	// 	fmt.Printf("%d: %s (%s)\n", i, tok.Value, tok.Type)
	// }

	diagram, err := interpreter.Parse(tokens)
	// fmt.Printf("Nodes: %d, Edges: %d\n", len(diagram.Nodes), len(diagram.Edges))
	if err != nil {
		fmt.Println("Parsing error:", err)
		return ""
	}

	svg := renderer.RenderSVG(diagram)
	base := strings.TrimSuffix(filepath.Base(path), ".diag")
	outPath := filepath.Join(outputDir, base+".svg")

	err = os.WriteFile(outPath, []byte(svg), 0644)
	if err != nil {
		fmt.Println("Could not save SVG:", err)
		return ""
	}
	// fmt.Println("Created:", outPath)
	return outPath

}

func renderAllDiagrams(diagramFiles []string) string {
	var output string
	for _, file := range diagramFiles {
		path := filepath.Join("example", file)
		svgPath := renderDiagToSVG(path)
		if svgPath != "" {
			output += fmt.Sprintf("Created: %s\n", svgPath)
		}
	}
	return output
}
