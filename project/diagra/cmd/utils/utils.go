package utils

import (
	"diagra/interpreter"
	"diagra/renderer"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var (
	RenderStart  time.Time
	CombinedTime int64
	mu           sync.Mutex
)

const (
	ExampleDir = "example"
	OutputDir  = "output"
)

// CheckError checks if an error occurred and prints it to the console.
func CheckError(err error) {
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

// ResetRenderStart resets the RenderStart time to the current time.
// it uses mutex because i was playing around with goroutines
func ResetRenderStart() {
	mu.Lock()
	defer mu.Unlock()
	RenderStart = time.Now()
}

// ResetCombinedTime resets the CombinedTime to 0.
func ResetCombinedTime() {
	mu.Lock()
	defer mu.Unlock()
	CombinedTime = 0
}

// RenderDiagToSVG reads a .diag file, parses it, and renders it to an SVG file.
// It creates an output directory if it doesn't exist and saves the SVG file there
func RenderDiagToSVG(path string) string {
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

// RenderAllDiagrams renders all diagrams in the given list of diagram files.
// It reads each file, processes it, and saves the output as SVG files.
func RenderAllDiagrams(diagramFiles []string) string {
	var output string
	for _, file := range diagramFiles {
		path := filepath.Join("example", file)
		svgPath := RenderDiagToSVG(path)
		CombinedTime += time.Since(RenderStart).Milliseconds()
		if svgPath != "" {
			output += fmt.Sprintf("Created: %s\n", svgPath)

		}
	}
	return output
}
