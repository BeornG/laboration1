package interpreter

import "fmt"

// Parser struct for parsing diagram definitions
type parser struct {
	tokens  []Token
	current int
}

// Parse starts parsing the tokens and returns a Diagram object
func Parse(tokens []Token) (Diagram, error) {
	p := &parser{tokens: tokens, current: 0}
	return p.parseDiagram()
}

// --- Internal help functions ---

// currentToken returns the current token being parsed
// If the current index is out of bounds, it returns an EOF token
func (p *parser) currentToken() Token {
	if p.current >= len(p.tokens) {
		return Token{Type: TOKEN_EOF}
	}
	return p.tokens[p.current]
}

// advance moves the current token index forward
func (p *parser) advance() {
	p.current++
}

// match checks if the current token matches the expected type
func (p *parser) match(typ TokenType) bool {
	if p.currentToken().Type == typ {
		p.advance()
		return true
	}
	return false
}

func (p *parser) parseDiagram() (Diagram, error) {
	var d Diagram

	// Expect: "diagram"
	if p.currentToken().Type != TOKEN_KEYWORD || p.currentToken().Value != "diagram" {
		return d, fmt.Errorf("expected 'diagram' keyword")
	}
	p.advance()

	// Expect: diagram type name
	if p.currentToken().Type != TOKEN_IDENTIFIER {
		return d, fmt.Errorf("expected diagram type name")
	}
	d.Name = p.currentToken().Value
	p.advance()

	if !allowedTypes[d.Name] {
		return d, fmt.Errorf("okänd diagramtyp: %s", d.Name)
	}

	// Optional attribute: (layout)
	if p.currentToken().Value == "(" {
		p.advance()
		for {
			if p.currentToken().Value == ")" || p.currentToken().Type == TOKEN_EOF {
				break
			}

			key := p.currentToken().Value
			p.advance()

			if p.currentToken().Value != "=" {
				return d, fmt.Errorf("expected '=' after attributename")
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

	// Expect: "{"
	// This is where the diagram content starts
	if !p.match(TOKEN_LBRACE) {
		return d, fmt.Errorf("expected '{' after diagram type")
	}

	for p.currentToken().Type != TOKEN_RBRACE && p.currentToken().Type != TOKEN_EOF {
		tok := p.currentToken()

		// --- Nodes ---
		if tok.Type == TOKEN_KEYWORD && tok.Value == "node" {
			p.advance()
			id := p.currentToken().Value
			p.advance()
			label := p.currentToken().Value
			p.advance()

			// Default values
			color := "#e0f7fa"     // light cyan
			textColor := "#004d40" // dark cyan
			shape := "rect"        // default shape: rectangle
			border := "#00796b"    // default border color: dark cyan

			if p.currentToken().Value == "(" {
				p.advance()
				for {
					if p.currentToken().Value == ")" || p.currentToken().Type == TOKEN_EOF {
						break
					}

					key := p.currentToken().Value
					p.advance()

					if p.currentToken().Value != "=" {
						return d, fmt.Errorf("expected '=' in node attribute")
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
				p.advance() // close ")"
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
				return d, fmt.Errorf("expected '->' after %s", from)
			}

			to := p.currentToken().Value
			p.advance()

			label := ""
			if p.currentToken().Type == TOKEN_STRING {
				label = p.currentToken().Value
				p.advance()
			}

			color := "#37474f" // default color: dark grey
			width := "2"       // default width: 2

			if p.currentToken().Value == "(" {
				p.advance()
				for {
					// fmt.Println("Edge attr loop at:", p.currentToken())
					if p.currentToken().Value == ")" || p.currentToken().Type == TOKEN_EOF {
						break
					}

					key := p.currentToken().Value
					p.advance()

					if p.currentToken().Value != "=" {
						return d, fmt.Errorf("expected '=' in edge attribute")
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
				p.advance() // close ")"
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
