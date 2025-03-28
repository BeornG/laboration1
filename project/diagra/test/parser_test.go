package interpreter_test

import (
	"diagra/interpreter"
	"testing"
)

func TestParser_SimpleDiagram(t *testing.T) {
	input := `
		diagram flowchart {
			node A "Start"
			node B "Bearbeta"
			A -> B "Går vidare"
		}
	`

	tokens := interpreter.Lex(input)
	diagram, err := interpreter.Parse(tokens)
	if err != nil {
		t.Fatalf("Fel vid tolkning: %v", err)
	}

	if len(diagram.Nodes) != 2 {
		t.Errorf("Förväntade 2 noder, fick %d", len(diagram.Nodes))
	}

	if len(diagram.Edges) != 1 {
		t.Errorf("Förväntade 1 kant, fick %d", len(diagram.Edges))
	}

	if diagram.Edges[0].Label != "Går vidare" {
		t.Errorf("Kantens label borde vara 'Går vidare', fick %s", diagram.Edges[0].Label)
	}
}
