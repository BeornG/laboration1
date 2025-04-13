package interpreter

import "fmt"

type Diagram struct {
	Name   string
	Layout string
	Nodes  []Node
	Edges  []Edge
}

type Node struct {
	ID    string
	Label string
}

type Edge struct {
	From  string
	To    string
	Label string
}

// Parserstruktur
type parser struct {
	tokens  []Token
	current int
}

// Global funktion för att starta parsning
func Parse(tokens []Token) (Diagram, error) {
	p := &parser{tokens: tokens, current: 0}
	return p.parseDiagram()
}

// --- Interna hjälpfunktioner ---

func (p *parser) currentToken() Token {
	if p.current >= len(p.tokens) {
		return Token{Type: TOKEN_EOF}
	}
	return p.tokens[p.current]
}

func (p *parser) advance() {
	p.current++
}

func (p *parser) match(typ TokenType) bool {
	if p.currentToken().Type == typ {
		p.advance()
		return true
	}
	return false
}

func (p *parser) parseDiagram() (Diagram, error) {
	var d Diagram

	// Förvänta: "diagram"
	if p.currentToken().Type != TOKEN_KEYWORD || p.currentToken().Value != "diagram" {
		return d, fmt.Errorf("förväntade 'diagram' som starttoken")
	}
	p.advance()

	// Förvänta: typnamn (t.ex. flowchart)
	if p.currentToken().Type != TOKEN_IDENTIFIER {
		return d, fmt.Errorf("förväntade diagramtyp efter 'diagram'")
	}
	d.Name = p.currentToken().Value
	p.advance()

	// Förvänta: {
	if !p.match(TOKEN_LBRACE) {
		return d, fmt.Errorf("förväntade '{' efter diagramtyp")
	}

	// Läs innehållet
	for p.currentToken().Type != TOKEN_RBRACE && p.currentToken().Type != TOKEN_EOF {
		tok := p.currentToken()

		// Noder
		if tok.Type == TOKEN_KEYWORD && tok.Value == "node" {
			p.advance()
			id := p.currentToken().Value
			p.advance()

			label := p.currentToken().Value
			p.advance()

			d.Nodes = append(d.Nodes, Node{ID: id, Label: label})
			continue
		}

		// Kanter
		if tok.Type == TOKEN_IDENTIFIER {
			from := tok.Value
			p.advance()

			if !p.match(TOKEN_ARROW) {
				return d, fmt.Errorf("förväntade '->' efter %s", from)
			}

			to := p.currentToken().Value
			p.advance()

			label := p.currentToken().Value
			p.advance()

			d.Edges = append(d.Edges, Edge{From: from, To: to, Label: label})
			continue
		}

		// Hoppa över oväntade tokens
		p.advance()
	}

	p.match(TOKEN_RBRACE)
	return d, nil
}
