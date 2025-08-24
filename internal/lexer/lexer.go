package lexer

// Lexer tokenizes JSON input
type Lexer struct {
	Input  string // Use string instead of []rune for better performance
	Pos    int
	Length int
}

// NewLexer creates a new lexer instance
func NewLexer(text string) *Lexer {
	return &Lexer{
		Input:  text,
		Pos:    0,
		Length: len(text),
	}
}

// Peek returns the current character without advancing the position
func (l *Lexer) Peek() byte {
	if l.Pos >= l.Length {
		return 0
	}
	return l.Input[l.Pos]
}

// Next returns the current character and advances the position
func (l *Lexer) Next() byte {
	if l.Pos >= l.Length {
		return 0
	}
	char := l.Input[l.Pos]
	l.Pos++
	return char
}

// SkipWhiteSpaces skips whitespace characters - optimized version
func (l *Lexer) SkipWhiteSpaces() {
	for l.Pos < l.Length {
		char := l.Input[l.Pos]
		if char == ' ' || char == '\t' || char == '\n' || char == '\r' {
			l.Pos++
		} else {
			break
		}
	}
}

// NextToken returns the next token from the input
func (l *Lexer) NextToken() Token {
	l.SkipWhiteSpaces()

	if l.Pos >= l.Length {
		return Token{EOF, ""}
	}

	char := l.Input[l.Pos]
	l.Pos++

	switch char {
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
	case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		l.Pos-- // Back up to include the digit/minus in parseNumber
		return l.parseNumber()
	case 't':
		return l.parseKeyword("true", True)
	case 'f':
		return l.parseKeyword("false", False)
	case 'n':
		return l.parseKeyword("null", Null)
	default:
		return Token{Invalid, string(char)}
	}
}

// parseString optimized to avoid string building
func (l *Lexer) parseString() Token {
	startPos := l.Pos

	for l.Pos < l.Length {
		char := l.Input[l.Pos]
		l.Pos++

		if char == '"' {
			return Token{String, l.Input[startPos : l.Pos-1]}
		}
		if char == '\\' && l.Pos < l.Length {
			l.Pos++ // Skip escaped character
		}
	}

	return Token{Invalid, "Unterminated string"}
}

// parseNumber optimized for better performance
func (l *Lexer) parseNumber() Token {
	startPos := l.Pos

	// Handle optional minus
	if l.Pos < l.Length && l.Input[l.Pos] == '-' {
		l.Pos++
	}

	// Parse digits
	if l.Pos >= l.Length || !isDigit(l.Input[l.Pos]) {
		return Token{Invalid, "Invalid number"}
	}

	// Handle zero or other digits
	if l.Input[l.Pos] == '0' {
		l.Pos++
	} else {
		for l.Pos < l.Length && isDigit(l.Input[l.Pos]) {
			l.Pos++
		}
	}

	// Handle decimal part
	if l.Pos < l.Length && l.Input[l.Pos] == '.' {
		l.Pos++
		if l.Pos >= l.Length || !isDigit(l.Input[l.Pos]) {
			return Token{Invalid, "Invalid number"}
		}
		for l.Pos < l.Length && isDigit(l.Input[l.Pos]) {
			l.Pos++
		}
	}

	// Handle exponent
	if l.Pos < l.Length && (l.Input[l.Pos] == 'e' || l.Input[l.Pos] == 'E') {
		l.Pos++
		if l.Pos < l.Length && (l.Input[l.Pos] == '+' || l.Input[l.Pos] == '-') {
			l.Pos++
		}
		if l.Pos >= l.Length || !isDigit(l.Input[l.Pos]) {
			return Token{Invalid, "Invalid number"}
		}
		for l.Pos < l.Length && isDigit(l.Input[l.Pos]) {
			l.Pos++
		}
	}

	return Token{Number, l.Input[startPos:l.Pos]}
}

// parseKeyword optimized for specific keywords
func (l *Lexer) parseKeyword(expected string, tokenType TokenType) Token {
	startPos := l.Pos - 1 // We already consumed the first character

	if l.Pos+len(expected)-1 > l.Length {
		return Token{Invalid, "Invalid keyword"}
	}

	// Check if the remaining characters match
	for i := 1; i < len(expected); i++ {
		if l.Input[l.Pos] != expected[i] {
			// Consume the rest of the invalid token
			for l.Pos < l.Length && isLetter(l.Input[l.Pos]) {
				l.Pos++
			}
			return Token{Invalid, l.Input[startPos:l.Pos]}
		}
		l.Pos++
	}

	return Token{tokenType, expected}
}

// Helper functions for better performance
func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

func isLetter(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}
