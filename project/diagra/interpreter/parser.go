package interpreter

import "fmt"

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

	if !allowedTypes[d.Name] {
		return d, fmt.Errorf("okänd diagramtyp: %s", d.Name)
	}

	// Valfria attribut (t.ex. layout)
	if p.currentToken().Value == "(" {
		p.advance()
		for {
			if p.currentToken().Value == ")" || p.currentToken().Type == TOKEN_EOF {
				break
			}

			key := p.currentToken().Value
			p.advance()

			if p.currentToken().Value != "=" {
				return d, fmt.Errorf("förväntade '=' efter attributnamn")
			}
			p.advance()

			value := p.currentToken().Value
			p.advance()

			if key == "layout" {
				d.Layout = value
			}

			if p.currentToken().Value == "," {
				p.advance()
			}
		}
		p.match(TOKEN_SYMBOL) // ")"
	}

	// Förvänta: {
	if !p.match(TOKEN_LBRACE) {
		return d, fmt.Errorf("förväntade '{' efter diagramtyp")
	}

	for p.currentToken().Type != TOKEN_RBRACE && p.currentToken().Type != TOKEN_EOF {
		tok := p.currentToken()

		// --- Noder ---
		if tok.Type == TOKEN_KEYWORD && tok.Value == "node" {
			p.advance()
			id := p.currentToken().Value
			p.advance()
			label := p.currentToken().Value
			p.advance()

			// default värden
			color := "#e0f7fa"
			textColor := "#004d40"
			shape := "rect"
			border := "#00796b"

			if p.currentToken().Value == "(" {
				p.advance()
				for {
					if p.currentToken().Value == ")" || p.currentToken().Type == TOKEN_EOF {
						break
					}

					key := p.currentToken().Value
					p.advance()

					if p.currentToken().Value != "=" {
						return d, fmt.Errorf("förväntade '=' i nod-attribut")
					}
					p.advance()

					value := p.currentToken().Value
					p.advance()

					switch key {
					case "color":
						color = value
					case "text":
						textColor = value
					case "shape":
						shape = value
					case "border":
						border = value
					}

					if p.currentToken().Value == "," {
						p.advance()
					}
				}
				p.advance() // stäng ")"
			}

			d.Nodes = append(d.Nodes, Node{
				ID:     id,
				Label:  label,
				Color:  color,
				Text:   textColor,
				Shape:  shape,
				Border: border,
			})
			continue
		}

		// --- Edges ---
		if tok.Type == TOKEN_IDENTIFIER {
			from := tok.Value
			p.advance()

			if !p.match(TOKEN_ARROW) {
				return d, fmt.Errorf("förväntade '->' efter %s", from)
			}

			to := p.currentToken().Value
			p.advance()

			label := ""
			if p.currentToken().Type == TOKEN_STRING {
				label = p.currentToken().Value
				p.advance()
			}

			color := "#37474f"
			width := "2"

			if p.currentToken().Value == "(" {
				p.advance()
				for {
					fmt.Println("Edge attr loop at:", p.currentToken())
					if p.currentToken().Value == ")" || p.currentToken().Type == TOKEN_EOF {
						break
					}

					key := p.currentToken().Value
					p.advance()

					if p.currentToken().Value != "=" {
						return d, fmt.Errorf("förväntade '=' i edge-attribut")
					}
					p.advance()

					value := p.currentToken().Value
					p.advance()

					switch key {
					case "color":
						color = value
					case "width":
						width = value
					}

					if p.currentToken().Value == "," {
						p.advance()
					}
				}
				p.advance() // stäng ")"
			}

			d.Edges = append(d.Edges, Edge{
				From:  from,
				To:    to,
				Label: label,
				Color: color,
				Width: width,
			})
			continue
		}

		p.advance()
	}

	p.match(TOKEN_RBRACE)
	return d, nil
}
