package parser

// JSONValue represents any JSON value
type JSONValue interface{}

// JSONObject represents a JSON object
type JSONObject map[string]JSONValue

// JSONArray represents a JSON array
type JSONArray []JSONValue

// JSONString represents a JSON string
type JSONString string

// JSONNumber represents a JSON number
type JSONNumber float64

// JSONBool represents a JSON boolean
type JSONBool bool

// JSONNull represents a JSON null value
type JSONNull struct{}

// ParseError represents a parsing error
type ParseError struct {
	Message string
	Pos     int
}

func (e *ParseError) Error() string {
	return e.Message
}
