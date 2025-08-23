package lexer

import "unicode"

// Lexer tokenizes JSON input
type Lexer struct {
	Input []rune
	Pos   int
}

// NewLexer creates a new lexer instance
func NewLexer(text string) *Lexer {
	return &Lexer{
		Input: []rune(text),
		Pos:   0,
	}
}

// Peek returns the current character without advancing the position
func (l *Lexer) Peek() rune {
	if l.Pos >= len(l.Input) {
		return 0
	}
	return l.Input[l.Pos]
}

// Next returns the current character and advances the position
func (l *Lexer) Next() rune {
	char := l.Peek()
	l.Pos++
	return char
}

// SkipWhiteSpaces skips whitespace characters
func (l *Lexer) SkipWhiteSpaces() {
	for unicode.IsSpace(l.Peek()) {
		l.Next()
	}
}

// NextToken returns the next token from the input
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
		return l.parseString()
	default:
		if unicode.IsDigit(char) || char == '-' {
			return l.parseNumber(char)
		}
		return l.parseKeyword(char)
	}
}

// parseString parses a JSON string token
func (l *Lexer) parseString() Token {
	startPos := l.Pos
	for {
		char := l.Next()
		if char == '"' {
			break
		}
		if char == 0 {
			panic("Unterminated string")
		}
		// Handle escape sequences
		if char == '\\' {
			l.Next() // Skip the escaped character
		}
	}
	return Token{String, string(l.Input[startPos : l.Pos-1])}
}

// parseNumber parses a JSON number token
func (l *Lexer) parseNumber(firstChar rune) Token {
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
			// Handle optional sign after exponent
			if l.Peek() == '+' || l.Peek() == '-' {
				l.Next()
			}
		} else {
			break
		}
	}

	return Token{Number, string(l.Input[startPos:l.Pos])}
}

// parseKeyword parses JSON keywords (true, false, null)
func (l *Lexer) parseKeyword(firstChar rune) Token {
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
		// Return an invalid token instead of panicking
		return Token{Invalid, word}
	}
}
