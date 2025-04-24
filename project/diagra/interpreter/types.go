package interpreter

type TokenType string

const (
	TOKEN_KEYWORD    TokenType = "KEYWORD"
	TOKEN_IDENTIFIER TokenType = "IDENTIFIER"
	TOKEN_STRING     TokenType = "STRING"
	TOKEN_ARROW      TokenType = "ARROW"
	TOKEN_LBRACE     TokenType = "LBRACE"
	TOKEN_RBRACE     TokenType = "RBRACE"
	TOKEN_SYMBOL     TokenType = "SYMBOL"
	TOKEN_EOF        TokenType = "EOF"
)

type Token struct {
	Type  TokenType
	Value string
}

type Diagram struct {
	Name   string
	Layout string
	Nodes  []Node
	Edges  []Edge
}

type Node struct {
	ID     string
	Label  string
	Color  string
	Text   string
	Shape  string
	Border string
}

type Edge struct {
	From  string
	To    string
	Label string
	Color string
	Width string
}

var allowedTypes = map[string]bool{
	"flowchart": true,
	"tree":      true,
}
