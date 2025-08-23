package main

import (
	"fmt"
	"log"

	"github.com/moabdelazem/parser/internal/lexer"
	"github.com/moabdelazem/parser/internal/parser"
)

/*
	Our goal is to write a JSON parser from scratch in Go.
	RFC 8259 will be our guide: https://datatracker.ietf.org/doc/html/rfc8259

	JSON can represent four primitive types (strings, numbers, booleans, and null)
	and two structured types (objects and arrays).
*/

func main() {
	fmt.Println("ok let's try again")

	// Example
	input := `{"name": "Alice", "age": 25, "isStudent": false, "grades": [95, 87, 92]}`

	fmt.Println("\n1. Tokenization:")
	lex := lexer.NewLexer(input)

	fmt.Printf("Input: %s\n\n", input)
	fmt.Println("Tokens:")

	lexForParsing := lexer.NewLexer(input)
	for {
		tok := lex.NextToken()
		fmt.Printf("%-10s: %s\n", tok.Type.String(), tok.Value)
		if tok.Type == lexer.EOF {
			break
		}
	}

	fmt.Println("\n2. Parsing:")
	p := parser.NewParser(lexForParsing)
	result, err := p.Parse()
	if err != nil {
		log.Fatalf("Parse error: %v", err)
	}

	fmt.Printf("Parsed JSON structure: %+v\n", result)

	fmt.Println("\n3. Accessing parsed data:")
	if obj, ok := result.(parser.JSONObject); ok {
		if name, exists := obj["name"]; exists {
			fmt.Printf("Name: %v\n", name)
		}
		if age, exists := obj["age"]; exists {
			fmt.Printf("Age: %v\n", age)
		}
		if grades, exists := obj["grades"]; exists {
			if gradeArray, ok := grades.(parser.JSONArray); ok {
				fmt.Printf("Grades: %v\n", gradeArray)
			}
		}
	}
}
