package main

import (
	"fmt"
	"log"
	"strings"

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
	fmt.Println("JSON Parser Examples")

	// Example 1: Student Object
	example1 := `{"name": "Alice", "age": 25, "isStudent": false, "grades": [95, 87, 92]}`
	runExample("Example 1: Student Object", example1)

	fmt.Println("\n" + strings.Repeat("=", 60))

	// Example 2: Complex Nested Structure
	example2 := `{
		"company": "TechCorp",
		"employees": [
			{
				"id": 1,
				"name": "John Doe",
				"position": "Software Engineer",
				"salary": 75000.50,
				"skills": ["Go", "JavaScript", "Python"],
				"remote": true,
				"manager": null
			},
			{
				"id": 2,
				"name": "Jane Smith",
				"position": "Product Manager",
				"salary": 85000,
				"skills": ["Leadership", "Analytics"],
				"remote": false,
				"manager": {
					"name": "Bob Wilson",
					"department": "Engineering"
				}
			}
		],
		"founded": 2015,
		"public": true
	}`
	runExample("Example 2: Complex Nested Structure", example2)
}

func runExample(title string, input string) {
	fmt.Printf("\n%s\n", title)

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
		// Access data based on the structure
		if title == "Example 1: Student Object" {
			accessStudentData(obj)
		} else {
			accessCompanyData(obj)
		}
	}
}

func accessStudentData(obj parser.JSONObject) {
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

func accessCompanyData(obj parser.JSONObject) {
	if company, exists := obj["company"]; exists {
		fmt.Printf("Company: %v\n", company)
	}
	if founded, exists := obj["founded"]; exists {
		fmt.Printf("Founded: %v\n", founded)
	}
	if employees, exists := obj["employees"]; exists {
		if empArray, ok := employees.(parser.JSONArray); ok {
			fmt.Printf("Number of employees: %d\n", len(empArray))

			// Access first employee details
			if len(empArray) > 0 {
				if firstEmp, ok := empArray[0].(parser.JSONObject); ok {
					if name, exists := firstEmp["name"]; exists {
						fmt.Printf("First employee name: %v\n", name)
					}
					if salary, exists := firstEmp["salary"]; exists {
						fmt.Printf("First employee salary: %v\n", salary)
					}
					if skills, exists := firstEmp["skills"]; exists {
						if skillArray, ok := skills.(parser.JSONArray); ok {
							fmt.Printf("First employee skills: %v\n", skillArray)
						}
					}
				}
			}
		}
	}
}
