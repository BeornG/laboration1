package interpreter

import (
	"unicode"
)

var keywords = map[string]bool{
	"diagram": true,
	"node":    true,
}

// Lex takes a string input and returns a slice of tokens.
// It identifies keywords, identifiers, numbers, strings, and symbols.
// It also handles whitespace and comments.
func Lex(input string) []Token {
	// fmt.Println("Lexing started")
	var tokens []Token
	runes := []rune(input)
	length := len(runes)

	i := 0
	for i < length {
		c := runes[i]
		// fmt.Printf("Lex at %d: %q\n", i, c)
		if unicode.IsSpace(c) {
			i++
			continue
		}

		// Identifier and keywords
		if unicode.IsLetter(c) {
			start := i
			for i < length && (unicode.IsLetter(runes[i]) || unicode.IsDigit(runes[i])) {
				i++
			}
			value := string(runes[start:i])
			if keywords[value] {
				tokens = append(tokens, Token{Type: TOKEN_KEYWORD, Value: value})
			} else {
				tokens = append(tokens, Token{Type: TOKEN_IDENTIFIER, Value: value})
			}
			continue
		}

		// Numbers example: width = 3
		if unicode.IsDigit(c) {
			start := i
			for i < length && unicode.IsDigit(runes[i]) {
				i++
			}
			tokens = append(tokens, Token{Type: TOKEN_IDENTIFIER, Value: string(runes[start:i])})
			continue
		}

		// Strings "..."
		if c == '"' {
			i++
			start := i
			for i < length && runes[i] != '"' {
				i++
			}
			value := string(runes[start:i])
			tokens = append(tokens, Token{Type: TOKEN_STRING, Value: value})
			i++ // hoppa över slut-quote
			continue
		}

		// Arrows ->
		if c == '-' && i+1 < length && runes[i+1] == '>' {
			tokens = append(tokens, Token{Type: TOKEN_ARROW, Value: "->"})
			i += 2
			continue
		}

		// Braces
		if c == '{' {
			tokens = append(tokens, Token{Type: TOKEN_LBRACE, Value: "{"})
			i++
			continue
		}
		if c == '}' {
			tokens = append(tokens, Token{Type: TOKEN_RBRACE, Value: "}"})
			i++
			continue
		}

		// Check if it is '=', '(', ')', eller ','.
		// If it is, create TOKEN_SYMBOL and add to token list.
		if c == '=' || c == '(' || c == ')' || c == ',' {
			tokens = append(tokens, Token{Type: TOKEN_SYMBOL, Value: string(c)})
			i++
			continue
		}

		// Unknown, skip
		i++
	}

	tokens = append(tokens, Token{Type: TOKEN_EOF, Value: ""})
	return tokens
}
