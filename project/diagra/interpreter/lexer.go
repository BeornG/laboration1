package interpreter

import (
	"fmt"
	"unicode"
)

var keywords = map[string]bool{
	"diagram": true,
	"node":    true,
}

// Lex tar en sträng och returnerar en lista av tokens
func Lex(input string) []Token {
	fmt.Println("Lexing started")
	var tokens []Token
	runes := []rune(input)
	length := len(runes)

	i := 0
	for i < length {
		c := runes[i]
		// fmt.Printf("Lex at %d: %q\n", i, c)
		// Hoppa över whitespace
		if unicode.IsSpace(c) {
			i++
			continue
		}

		// Identifierare eller nyckelord
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

		// Siffror (för t.ex. width=3)
		if unicode.IsDigit(c) {
			start := i
			for i < length && unicode.IsDigit(runes[i]) {
				i++
			}
			tokens = append(tokens, Token{Type: TOKEN_IDENTIFIER, Value: string(runes[start:i])})
			continue
		}

		// Strängar "..."
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

		// Pilar ->
		if c == '-' && i+1 < length && runes[i+1] == '>' {
			tokens = append(tokens, Token{Type: TOKEN_ARROW, Value: "->"})
			i += 2
			continue
		}

		// Klammrar
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

		// Kontrollera om tecknet är ett av de specifika symbolerna '=', '(', ')', eller ','.
		// Om det är det, skapa en TOKEN_SYMBOL och lägg till den i tokens-listan.
		if c == '=' || c == '(' || c == ')' || c == ',' {
			tokens = append(tokens, Token{Type: TOKEN_SYMBOL, Value: string(c)})
			i++
			continue
		}

		// Okänt tecken – hoppa över
		i++
	}

	tokens = append(tokens, Token{Type: TOKEN_EOF, Value: ""})
	return tokens
}
