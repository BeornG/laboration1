package tui

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
)

// RunTUI starts the TUI application
// It loads the .diag files from the example directory and initializes the model
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
