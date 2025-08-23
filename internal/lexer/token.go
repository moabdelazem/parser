package lexer

// TokenType represents the type of a token in JSON
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
	Invalid                   // invalid token
	EOF                       // end of input
)

// Token represents a token with its type and value
type Token struct {
	Type  TokenType
	Value string
}

// String returns a string representation of the token type
func (t TokenType) String() string {
	switch t {
	case LBrace:
		return "LBrace"
	case RBrace:
		return "RBrace"
	case LBracket:
		return "LBracket"
	case RBracket:
		return "RBracket"
	case Colon:
		return "Colon"
	case Comma:
		return "Comma"
	case String:
		return "String"
	case Number:
		return "Number"
	case True:
		return "True"
	case False:
		return "False"
	case Null:
		return "Null"
	case Invalid:
		return "Invalid"
	case EOF:
		return "EOF"
	default:
		return "Unknown"
	}
}
