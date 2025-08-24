package parser

import (
	"strings"
	"testing"

	"github.com/moabdelazem/parser/internal/lexer"
)

func BenchmarkSmallObject(b *testing.B) {
	input := `{"name": "Alice", "age": 25, "isStudent": false}`
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		lex := lexer.NewLexer(input)
		p := NewParser(lex)
		_, err := p.Parse()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkLargeArray(b *testing.B) {
	// Create a large array
	elements := make([]string, 1000)
	for i := range elements {
		elements[i] = `{"id": ` + string(rune(i)) + `, "name": "item` + string(rune(i)) + `"}`
	}
	input := "[" + strings.Join(elements, ",") + "]"

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		lex := lexer.NewLexer(input)
		p := NewParser(lex)
		_, err := p.Parse()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkNestedObjects(b *testing.B) {
	input := `{"level1": {"level2": {"level3": {"value": 42}}}}`
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		lex := lexer.NewLexer(input)
		p := NewParser(lex)
		_, err := p.Parse()
		if err != nil {
			b.Fatal(err)
		}
	}
}
