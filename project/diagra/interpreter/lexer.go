package interpreter

import (
	"unicode"
)

type TokenType string

const (
	TOKEN_KEYWORD    TokenType = "KEYWORD"
	TOKEN_IDENTIFIER TokenType = "IDENTIFIER"
	TOKEN_STRING     TokenType = "STRING"
	TOKEN_ARROW      TokenType = "ARROW"
	TOKEN_LBRACE     TokenType = "LBRACE"
	TOKEN_RBRACE     TokenType = "RBRACE"
	TOKEN_EOF        TokenType = "EOF"
)

type Token struct {
	Type  TokenType
	Value string
}

var keywords = map[string]bool{
	"diagram": true,
	"node":    true,
}

// Lex tar en sträng och returnerar en lista av tokens
func Lex(input string) []Token {
	var tokens []Token
	runes := []rune(input)
	length := len(runes)

	i := 0
	for i < length {
		c := runes[i]

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

		// Okänt tecken – hoppa över
		i++
	}

	tokens = append(tokens, Token{Type: TOKEN_EOF, Value: ""})
	return tokens
}
