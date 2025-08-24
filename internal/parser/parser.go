package parser

import (
	"strconv"
	"sync"

	"github.com/moabdelazem/parser/internal/lexer"
)

var (
	// Object and array pools to reduce allocations
	objectPool = sync.Pool{
		New: func() interface{} {
			return make(JSONObject)
		},
	}

	arrayPool = sync.Pool{
		New: func() interface{} {
			return make(JSONArray, 0, 8) // Pre-allocate with capacity
		},
	}
)

type Parser struct {
	lexer        *lexer.Lexer
	currentToken lexer.Token
	peekToken    lexer.Token // Add peek token for lookahead
}

func NewParser(l *lexer.Lexer) *Parser {
	p := &Parser{
		lexer: l,
	}
	p.advance() // Load current token
	p.advance() // Load peek token
	return p
}

func (p *Parser) advance() {
	p.currentToken = p.peekToken
	p.peekToken = p.lexer.NextToken()
}

func (p *Parser) expectToken(tokenType lexer.TokenType) error {
	if p.currentToken.Type != tokenType {
		return &ParseError{
			Message: "Expected " + tokenType.String() + " but got " + p.currentToken.Type.String(),
			Pos:     p.lexer.Pos,
		}
	}
	p.advance()
	return nil
}

func (p *Parser) Parse() (JSONValue, error) {
	value, err := p.parseValue()
	if err != nil {
		return nil, err
	}

	if p.currentToken.Type != lexer.EOF {
		return nil, &ParseError{
			Message: "Expected end of input but got " + p.currentToken.Type.String(),
			Pos:     p.lexer.Pos,
		}
	}

	return value, nil
}

func (p *Parser) parseValue() (JSONValue, error) {
	switch p.currentToken.Type {
	case lexer.LBrace:
		return p.parseObject()
	case lexer.LBracket:
		return p.parseArray()
	case lexer.String:
		value := JSONString(p.currentToken.Value)
		p.advance()
		return value, nil
	case lexer.Number:
		// Optimize number parsing - avoid string conversion when possible
		return p.parseNumberValue()
	case lexer.True:
		p.advance()
		return JSONBool(true), nil
	case lexer.False:
		p.advance()
		return JSONBool(false), nil
	case lexer.Null:
		p.advance()
		return JSONNull{}, nil
	default:
		return nil, &ParseError{
			Message: "Unexpected token: " + p.currentToken.Type.String(),
			Pos:     p.lexer.Pos,
		}
	}
}

func (p *Parser) parseNumberValue() (JSONValue, error) {
	numStr := p.currentToken.Value
	p.advance()

	// Try to parse as int first for better performance
	if val, err := strconv.Atoi(numStr); err == nil {
		return JSONNumber(float64(val)), nil
	}

	// Fall back to float parsing
	if val, err := strconv.ParseFloat(numStr, 64); err == nil {
		return JSONNumber(val), nil
	}

	return nil, &ParseError{
		Message: "Invalid number: " + numStr,
		Pos:     p.lexer.Pos,
	}
}

func (p *Parser) parseObject() (JSONObject, error) {
	obj := objectPool.Get().(JSONObject)
	// Clear the map
	for k := range obj {
		delete(obj, k)
	}

	if err := p.expectToken(lexer.LBrace); err != nil {
		objectPool.Put(obj)
		return nil, err
	}

	if p.currentToken.Type == lexer.RBrace {
		p.advance()
		return obj, nil
	}

	for {
		if p.currentToken.Type != lexer.String {
			objectPool.Put(obj)
			return nil, &ParseError{
				Message: "Expected string key in object",
				Pos:     p.lexer.Pos,
			}
		}
		key := p.currentToken.Value
		p.advance()

		if err := p.expectToken(lexer.Colon); err != nil {
			objectPool.Put(obj)
			return nil, err
		}

		value, err := p.parseValue()
		if err != nil {
			objectPool.Put(obj)
			return nil, err
		}
		obj[key] = value

		if p.currentToken.Type == lexer.RBrace {
			p.advance()
			break
		}
		if p.currentToken.Type != lexer.Comma {
			objectPool.Put(obj)
			return nil, &ParseError{
				Message: "Expected ',' or '}' in object",
				Pos:     p.lexer.Pos,
			}
		}
		p.advance()
	}

	return obj, nil
}

func (p *Parser) parseArray() (JSONArray, error) {
	arr := arrayPool.Get().(JSONArray)
	arr = arr[:0] // Reset length but keep capacity

	if err := p.expectToken(lexer.LBracket); err != nil {
		arrayPool.Put(&arr)
		return nil, err
	}

	if p.currentToken.Type == lexer.RBracket {
		p.advance()
		return arr, nil
	}

	for {
		value, err := p.parseValue()
		if err != nil {
			arrayPool.Put(&arr)
			return nil, err
		}
		arr = append(arr, value)

		if p.currentToken.Type == lexer.RBracket {
			p.advance()
			break
		}
		if p.currentToken.Type != lexer.Comma {
			arrayPool.Put(&arr)
			return nil, &ParseError{
				Message: "Expected ',' or ']' in array",
				Pos:     p.lexer.Pos,
			}
		}
		p.advance()
	}

	return arr, nil
}
