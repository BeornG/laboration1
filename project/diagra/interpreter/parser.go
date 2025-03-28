package interpreter

type Diagram struct {
	Name  string
	Nodes []Node
	Edges []Edge
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

	// diagram flowchart {
	if !p.match(TOKEN_KEYWORD) || p.currentToken().Value != "flowchart" {
		return d, nil
	}
	p.advance()

	if !p.match(TOKEN_LBRACE) {
		return d, nil
	}

	// Loopa tills vi hittar }
	for p.currentToken().Type != TOKEN_RBRACE && p.currentToken().Type != TOKEN_EOF {
		tok := p.currentToken()

		if tok.Type == TOKEN_KEYWORD && tok.Value == "node" {
			p.advance()
			id := p.currentToken().Value
			p.advance()
			label := p.currentToken().Value
			p.advance()

			d.Nodes = append(d.Nodes, Node{
				ID:    id,
				Label: label,
			})
			continue
		}

		if tok.Type == TOKEN_IDENTIFIER {
			from := tok.Value
			p.advance() // identifier
			if !p.match(TOKEN_ARROW) {
				return d, nil
			}
			to := p.currentToken().Value
			p.advance()
			label := p.currentToken().Value
			p.advance()

			d.Edges = append(d.Edges, Edge{
				From:  from,
				To:    to,
				Label: label,
			})
			continue
		}

		// Hoppa över okända tokens
		p.advance()
	}

	// }
	p.match(TOKEN_RBRACE)

	return d, nil
}
