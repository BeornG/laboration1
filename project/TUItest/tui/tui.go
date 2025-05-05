package tui

import (
	"fmt"
	"os"

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
	// var files []string
	// err := filepath.WalkDir(path, func(p string, d fs.DirEntry, e error) error {
	// 	if e != nil {
	// 		return e
	// 	}
	// 	if !d.IsDir() && filepath.Ext(p) == ".diag" {
	// 		files = append(files, filepath.Base(p))
	// 	}
	// 	return nil
	// })
	files := []string{
		"example1.diag",
		"example2.diag",
		"example3.diag",
		"exampleTree1.diag",
	}
	return files, nil
}

// Placeholder
// func renderDiagToSVG(diagFile string) string {
// 	// Simulering
// 	return fmt.Sprintf("Rendered %s.svg successfully!", diagFile)
// }
