package main

import (
	"diagra/interpreter"
	"diagra/renderer"
	"fmt"
	"os"
)

func main() {
	src, err := os.ReadFile("example/example2.diag")
	if err != nil {
		fmt.Println("Kunde inte läsa fil:", err)
		return
	}
	fmt.Println("Input:\n", string(src))
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
	err = os.WriteFile("output.svg", []byte(svg), 0644)
	if err != nil {
		fmt.Println("Kunde inte skriva fil:", err)
	}
}
