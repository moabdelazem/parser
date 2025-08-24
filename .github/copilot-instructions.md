# JSON Parser Copilot Instructions

## Architecture Overview

This is a from-scratch JSON parser following RFC 8259, built with a **two-phase design**: lexical analysis (tokenization) → syntactic analysis (parsing). The codebase demonstrates performance-optimized parsing techniques with memory pooling and zero-copy string operations.

### Key Components

- **`internal/lexer/`**: Tokenizes JSON input into structured tokens (LBrace, String, Number, etc.)
- **`internal/parser/`**: Recursive descent parser that builds AST from token stream
- **`cmd/parser/main.go`**: Demonstration CLI showing tokenization → parsing → data access workflow

## Critical Design Patterns

### 1. **Two-Phase Processing Pattern**

Always create separate lexer and parser instances:

```go
lex := lexer.NewLexer(input)
p := parser.NewParser(lex)
result, err := p.Parse()
```

### 2. **Performance-First Architecture**

- **Object Pooling**: Parser uses `sync.Pool` for JSONObject/JSONArray recycling
- **Zero-Copy Strings**: Lexer uses string slicing (`l.Input[startPos:l.Pos]`) instead of byte building
- **Lookahead Tokens**: Parser maintains `currentToken` + `peekToken` for efficient parsing decisions
- **Pre-allocated Capacities**: Arrays start with capacity 8 to reduce reallocations

### 3. **Type System Pattern**

JSON values use Go's type system with custom types:

```go
JSONObject = map[string]JSONValue  // Not map[string]interface{}
JSONArray = []JSONValue
JSONString = string (not string)   // Distinct types for type safety
```

### 4. **Error Handling Pattern**

- **Position-Aware Errors**: `ParseError` includes lexer position for debugging
- **Resource Cleanup**: Always `objectPool.Put(obj)` on parse errors
- **Fail-Fast Validation**: Lexer validates number formats during tokenization

## Development Workflows

### Testing Strategy

```bash
# Run specific package tests
go test ./internal/lexer -v
go test ./internal/parser -v

# Run performance benchmarks
go test ./internal/parser -bench=. -benchmem
```

### Performance Testing Pattern

Benchmarks focus on three scenarios:

- **Small objects** (simple key-value pairs)
- **Large arrays** (1000+ elements)
- **Nested objects** (deep hierarchy)

All benchmarks use `b.ReportAllocs()` to track memory allocations.

### Build & Run

```bash
# Standard execution
go run cmd/parser/main.go

# Build optimized binary
go build -o parser cmd/parser/main.go
```

## Coding Conventions

### 1. **Token Processing**

- Always check for `lexer.EOF` when iterating tokens
- Use `advance()` method consistently, never directly modify currentToken
- Handle `Invalid` token type in all switch statements

### 2. **Memory Management**

- Get objects from pools: `obj := objectPool.Get().(JSONObject)`
- Clear maps before reuse: `for k := range obj { delete(obj, k) }`
- Return objects to pools on ALL code paths (success + error)

### 3. **Number Parsing Optimization**

Parser attempts `strconv.Atoi()` first (faster for integers), falls back to `strconv.ParseFloat()`. Always store as `JSONNumber(float64)` for consistency.

### 4. **Test Structure**

- **Table-driven tests** for multiple inputs (see `TestLexerNumberFormats`)
- **Separate error cases** in dedicated test functions
- **Type assertions** with descriptive error messages: `if obj, ok := result.(JSONObject); !ok { t.Fatalf("Expected JSONObject, got %T", result) }`

## Integration Points

### Lexer ↔ Parser Interface

Parser depends on lexer's token stream. Key methods:

- `lexer.NextToken()` - primary interface, returns `Token{Type, Value}`
- Lexer maintains internal position (`l.Pos`) for error reporting

### Data Flow

1. **Input** → Lexer (string tokenization)
2. **Tokens** → Parser (AST construction)
3. **JSONValue** → Application (type-safe data access)

## Performance Considerations

### When to Use Pooling

Object/array pools are effective for:

- Parsing multiple JSON documents
- Large nested structures
- High-frequency parsing scenarios

### Streaming Support (Future Enhancement)

Current design assumes full input in memory. For file streaming, extend lexer with `io.Reader` interface while maintaining zero-copy principles where possible.

## Common Pitfalls

1. **Forgetting EOF checks** in token iteration loops
2. **Missing pool returns** in error paths (causes memory leaks)
3. **Type assertion without error checking** when accessing parsed data
4. **Creating new lexer for every token** (creates separate instances, breaks parsing)

## RFC 8259 Compliance Areas

- Whitespace handling (space, tab, newline, carriage return)
- String escape sequences (\", \\, \/, \b, \f, \n, \r, \t)
- Number formats (integers, decimals, scientific notation)
- Structural validation (proper nesting, no trailing commas)
