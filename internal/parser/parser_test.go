package parser

import (
	"testing"

	"github.com/moabdelazem/parser/internal/lexer"
)

func TestParseSimpleObject(t *testing.T) {
	input := `{"name": "Alice", "age": 25}`
	lex := lexer.NewLexer(input)
	parser := NewParser(lex)

	result, err := parser.Parse()
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	obj, ok := result.(JSONObject)
	if !ok {
		t.Fatalf("Expected JSONObject, got %T", result)
	}

	if len(obj) != 2 {
		t.Errorf("Expected 2 keys, got %d", len(obj))
	}

	name, exists := obj["name"]
	if !exists {
		t.Error("Key 'name' not found")
	}
	if nameStr, ok := name.(JSONString); !ok || string(nameStr) != "Alice" {
		t.Errorf("Expected name 'Alice', got %v", name)
	}

	age, exists := obj["age"]
	if !exists {
		t.Error("Key 'age' not found")
	}
	if ageNum, ok := age.(JSONNumber); !ok || float64(ageNum) != 25 {
		t.Errorf("Expected age 25, got %v", age)
	}
}

func TestParseArray(t *testing.T) {
	input := `[1, 2, 3]`
	lex := lexer.NewLexer(input)
	parser := NewParser(lex)

	result, err := parser.Parse()
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	arr, ok := result.(JSONArray)
	if !ok {
		t.Fatalf("Expected JSONArray, got %T", result)
	}

	if len(arr) != 3 {
		t.Errorf("Expected 3 elements, got %d", len(arr))
	}

	for i, expected := range []float64{1, 2, 3} {
		if num, ok := arr[i].(JSONNumber); !ok || float64(num) != expected {
			t.Errorf("Element %d: expected %v, got %v", i, expected, arr[i])
		}
	}
}

func TestParseComplexObject(t *testing.T) {
	input := `{"name": "Alice", "age": 25, "isStudent": false, "grades": [95, 87, 92], "address": null}`
	lex := lexer.NewLexer(input)
	parser := NewParser(lex)

	result, err := parser.Parse()
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	obj, ok := result.(JSONObject)
	if !ok {
		t.Fatalf("Expected JSONObject, got %T", result)
	}

	// Test boolean value
	isStudent, exists := obj["isStudent"]
	if !exists {
		t.Error("Key 'isStudent' not found")
	}
	if boolVal, ok := isStudent.(JSONBool); !ok || bool(boolVal) != false {
		t.Errorf("Expected isStudent false, got %v", isStudent)
	}

	// Test null value
	address, exists := obj["address"]
	if !exists {
		t.Error("Key 'address' not found")
	}
	if _, ok := address.(JSONNull); !ok {
		t.Errorf("Expected JSONNull, got %T", address)
	}

	// Test nested array
	grades, exists := obj["grades"]
	if !exists {
		t.Error("Key 'grades' not found")
	}
	if gradesArr, ok := grades.(JSONArray); !ok {
		t.Errorf("Expected JSONArray for grades, got %T", grades)
	} else if len(gradesArr) != 3 {
		t.Errorf("Expected 3 grades, got %d", len(gradesArr))
	}
}

func TestParseEmptyObjects(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"{}", "empty object"},
		{"[]", "empty array"},
	}

	for _, test := range tests {
		lex := lexer.NewLexer(test.input)
		parser := NewParser(lex)

		result, err := parser.Parse()
		if err != nil {
			t.Errorf("Parse error for %s: %v", test.expected, err)
			continue
		}

		switch test.input {
		case "{}":
			if obj, ok := result.(JSONObject); !ok || len(obj) != 0 {
				t.Errorf("Expected empty JSONObject, got %v", result)
			}
		case "[]":
			if arr, ok := result.(JSONArray); !ok || len(arr) != 0 {
				t.Errorf("Expected empty JSONArray, got %v", result)
			}
		}
	}
}

func TestParseError(t *testing.T) {
	invalidInputs := []string{
		`{"name": }`,         // missing value
		`{"name" "Alice"}`,   // missing colon
		`{"name": "Alice",}`, // trailing comma
		`{name: "Alice"}`,    // unquoted key
		`[1, 2,]`,            // trailing comma in array
	}

	for _, input := range invalidInputs {
		lex := lexer.NewLexer(input)
		parser := NewParser(lex)

		_, err := parser.Parse()
		if err == nil {
			t.Errorf("Expected error for invalid input: %s", input)
		}
	}
}
