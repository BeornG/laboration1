package cli

import (
	"diagra/cmd/utils"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

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
		fullPath := filepath.Join(utils.ExampleDir, args[1])
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			fmt.Println("File does not exist:", fullPath)
			return
		}
		renderCmd(fullPath)
		return
	case "render-all":
		renderAllCmd()
		utils.ResetCombinedTime()
		return
	case "-h", "--help", "help":
		helpCmd()
		return
	default:
		fmt.Println("Run -h, --help or help for usage")
		return
	}

}

// renderAllCmd renders all diagrams in the example directory.
// It reads all .diag files, processes them, and saves the output as SVG files.
func renderAllCmd() {
	path := utils.ExampleDir
	utils.ResetRenderStart()

	diagFiles, err := os.ReadDir(path)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return
	}
	var files []string
	for _, file := range diagFiles {
		if !strings.HasSuffix(file.Name(), ".diag") {
			continue
		}
		files = append(files, file.Name())
	}
	utils.RenderAllDiagrams(files)
	fmt.Println("All diagrams rendered to SVG in", utils.OutputDir)
	fmt.Printf("Total time: %d ms\n", utils.CombinedTime)
}

func renderCmd(filename string) {
	utils.ResetRenderStart()
	utils.RenderDiagToSVG(filename)
	fmt.Println("Rendering finished for", filename)
	timeTaken := time.Since(utils.RenderStart).Milliseconds()
	fmt.Printf("Total time: %d ms\n", timeTaken)
}

// helpCmd prints the help message for the CLI application.
// It shows the available commands and their usage.
func helpCmd() {
	fmt.Println("Usage: diagra [command]")
	fmt.Println("Commands:")
	fmt.Println("  render <file>		Render a diagram from a .diag file")
	fmt.Println("  render-all		Render all diagrams in the example directory")
	fmt.Println("  -h, --help, help     	Show this help message")
}
