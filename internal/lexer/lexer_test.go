package lexer

import (
	"testing"
)

func TestLexerBasicTokens(t *testing.T) {
	input := `{"name": "Alice", "age": 25}`
	lexer := NewLexer(input)

	expectedTokens := []struct {
		tokenType TokenType
		value     string
	}{
		{LBrace, "{"},
		{String, "name"},
		{Colon, ":"},
		{String, "Alice"},
		{Comma, ","},
		{String, "age"},
		{Colon, ":"},
		{Number, "25"},
		{RBrace, "}"},
		{EOF, ""},
	}

	for i, expected := range expectedTokens {
		token := lexer.NextToken()
		if token.Type != expected.tokenType {
			t.Errorf("Token %d: expected type %v, got %v", i, expected.tokenType, token.Type)
		}
		if token.Value != expected.value {
			t.Errorf("Token %d: expected value %q, got %q", i, expected.value, token.Value)
		}
	}
}

func TestLexerArray(t *testing.T) {
	input := `[1, 2.5, true, false, null]`
	lexer := NewLexer(input)

	expectedTokens := []struct {
		tokenType TokenType
		value     string
	}{
		{LBracket, "["},
		{Number, "1"},
		{Comma, ","},
		{Number, "2.5"},
		{Comma, ","},
		{True, "true"},
		{Comma, ","},
		{False, "false"},
		{Comma, ","},
		{Null, "null"},
		{RBracket, "]"},
		{EOF, ""},
	}

	for i, expected := range expectedTokens {
		token := lexer.NextToken()
		if token.Type != expected.tokenType {
			t.Errorf("Token %d: expected type %v, got %v", i, expected.tokenType, token.Type)
		}
		if token.Value != expected.value {
			t.Errorf("Token %d: expected value %q, got %q", i, expected.value, token.Value)
		}
	}
}

func TestLexerNumberFormats(t *testing.T) {
	inputs := []struct {
		input    string
		expected string
	}{
		{"123", "123"},
		{"-456", "-456"},
		{"3.14", "3.14"},
		{"-0.5", "-0.5"},
		{"1e10", "1e10"},
		{"2.5e-3", "2.5e-3"},
		{"1.23E+4", "1.23E+4"},
	}

	for _, test := range inputs {
		lexer := NewLexer(test.input)
		token := lexer.NextToken()

		if token.Type != Number {
			t.Errorf("Input %q: expected Number token, got %v", test.input, token.Type)
		}
		if token.Value != test.expected {
			t.Errorf("Input %q: expected value %q, got %q", test.input, test.expected, token.Value)
		}
	}
}
