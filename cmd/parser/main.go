package main

import (
	"fmt"
	"unicode"
)

/*
	Our goal is to write a JSON parser from scratch in Go.
	RFC 8259 will be our guide: https://datatracker.ietf.org/doc/html/rfc8259

	JSON can represent four primitive types (strings, numbers, booleans, and null)
	and two structured types (objects and arrays).
*/

type TokenType int

const (
	LBrace   TokenType = iota // {
	RBrace                    // }
	LBracket                  // [
	RBracket                  // ]
	Colon                     // :
	Comma                     // ,
	String                    // "text"
	Number                    // 123, 3.14
	True                      // true
	False                     // false
	Null                      // null
	EOF                       // end of input
)

type Token struct {
	Type  TokenType
	Value string
}

type Lexer struct {
	Input []rune
	Pos   int
}

func NewLexer(text string) *Lexer {
	return &Lexer{
		Input: []rune(text),
		Pos:   0,
	}
}

func (l *Lexer) Peek() rune {
	if l.Pos >= len(l.Input) {
		return 0
	}
	return l.Input[l.Pos]
}

func (l *Lexer) Next() rune {
	char := l.Peek()
	l.Pos++
	return char
}

func (l *Lexer) SkipWhiteSpaces() {
	for unicode.IsSpace(l.Peek()) {
		l.Next()
	}
}

func (l *Lexer) NextToken() Token {
	l.SkipWhiteSpaces()
	char := l.Next()

	switch char {
	case 0:
		return Token{EOF, ""}
	case '{':
		return Token{LBrace, "{"}
	case '}':
		return Token{RBrace, "}"}
	case '[':
		return Token{LBracket, "["}
	case ']':
		return Token{RBracket, "]"}
	case ':':
		return Token{Colon, ":"}
	case ',':
		return Token{Comma, ","}
	case '"':
		startPos := l.Pos
		for {
			char := l.Next()
			if char == '"' {
				break
			}
		}
		return Token{String, string(l.Input[startPos : l.Pos-1])}
	default:
		if unicode.IsDigit(char) || char == '-' {
			startPos := l.Pos - 1
			dotSeen := false
			expSeen := false
			for {
				c := l.Peek()
				if unicode.IsDigit(c) {
					l.Next()
				} else if c == '.' && !dotSeen && !expSeen {
					dotSeen = true
					l.Next()
				} else if (c == 'e' || c == 'E') && !expSeen {
					expSeen = true
					l.Next()
					// handle optional sign after exponent
					if l.Peek() == '+' || l.Peek() == '-' {
						l.Next()
					}
				} else {
					break
				}
			}

			return Token{Number, string(l.Input[startPos:l.Pos])}
		}

		startPos := l.Pos - 1
		for unicode.IsLetter(l.Peek()) {
			l.Next()
		}

		word := string(l.Input[startPos:l.Pos])
		switch word {
		case "true":
			return Token{True, word}
		case "false":
			return Token{False, word}
		case "null":
			return Token{Null, word}
		default:
			panic("Unexpected token: " + word)
		}
	}
}

func main() {
	fmt.Println("so we are writing a json parser from scratch")

	input := `{"name": "Alice", "age": 25, "isStudent": false}`
	lexer := NewLexer(input)

	for {
		tok := lexer.NextToken()
		fmt.Printf("Token: %-8v Value: %v\n", tok.Type, tok.Value)
		if tok.Type == EOF {
			break
		}
	}
}
