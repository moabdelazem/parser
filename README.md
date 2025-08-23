# JSON Parser from Scratch

A Simple JSON parser built from scratch in Go

## Features

- Tokenization of JSON input
- Parsing JSON objects `{}`
- Parsing JSON arrays `[]`
- Support for strings, numbers, booleans, and null values
- Proper error handling with descriptive messages

## Usage

```bash
# Run the parser
go run cmd/parser/main.go

# Build the binary
go build -o parser cmd/parser/main.go
```

## Example

```go
package main

import (
    "fmt"
    "log"

    "github.com/moabdelazem/parser/internal/lexer"
    "github.com/moabdelazem/parser/internal/parser"
)

func main() {
    input := `{"name": "Alice", "age": 25, "grades": [95, 87, 92]}`

    // Create lexer and parser
    lex := lexer.NewLexer(input)
    p := parser.NewParser(lex)

    // Parse JSON
    result, err := p.Parse()
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Parsed: %+v\n", result)
}
```

## JSON Types Supported

| JSON Type | Go Type                                    |
| --------- | ------------------------------------------ |
| Object    | `parser.JSONObject` (map[string]JSONValue) |
| Array     | `parser.JSONArray` ([]JSONValue)           |
| String    | `parser.JSONString` (string)               |
| Number    | `parser.JSONNumber` (float64)              |
| Boolean   | `parser.JSONBool` (bool)                   |
| Null      | `parser.JSONNull` (struct{})               |

## RFC 8259 Compliance

This parser follows the JSON specification as defined in RFC 8259:

- Proper handling of whitespace
- Support for escape sequences in strings
- Correct number parsing (integers, floats, scientific notation)
- Validation of JSON structure
