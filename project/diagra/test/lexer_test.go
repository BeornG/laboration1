package interpreter_test

import (
	"diagra/interpreter"
	"testing"
)

func TestLexer_TokensFromSimpleDiagram(t *testing.T) {
	input := `
		diagram flowchart {
			node A "Start"
			node B "Bearbeta"
			A -> B "Går vidare"
		}
	`

	expected := []interpreter.Token{
		{Type: interpreter.TOKEN_KEYWORD, Value: "diagram"},
		{Type: interpreter.TOKEN_IDENTIFIER, Value: "flowchart"},
		{Type: interpreter.TOKEN_LBRACE, Value: "{"},

		{Type: interpreter.TOKEN_KEYWORD, Value: "node"},
		{Type: interpreter.TOKEN_IDENTIFIER, Value: "A"},
		{Type: interpreter.TOKEN_STRING, Value: "Start"},

		{Type: interpreter.TOKEN_KEYWORD, Value: "node"},
		{Type: interpreter.TOKEN_IDENTIFIER, Value: "B"},
		{Type: interpreter.TOKEN_STRING, Value: "Bearbeta"},

		{Type: interpreter.TOKEN_IDENTIFIER, Value: "A"},
		{Type: interpreter.TOKEN_ARROW, Value: "->"},
		{Type: interpreter.TOKEN_IDENTIFIER, Value: "B"},
		{Type: interpreter.TOKEN_STRING, Value: "Går vidare"},

		{Type: interpreter.TOKEN_RBRACE, Value: "}"},
		{Type: interpreter.TOKEN_EOF, Value: ""},
	}

	tokens := interpreter.Lex(input)

	if len(tokens) != len(expected) {
		t.Fatalf("Förväntade %d tokens, fick %d", len(expected), len(tokens))
	}

	for i, token := range tokens {
		if token.Type != expected[i].Type || token.Value != expected[i].Value {
			t.Errorf("Token %d felaktig: förväntade (%s, %q), fick (%s, %q)",
				i,
				expected[i].Type, expected[i].Value,
				token.Type, token.Value)
		}
	}
}
