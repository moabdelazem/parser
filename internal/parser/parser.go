package parser

import (
	"strconv"

	"github.com/moabdelazem/parser/internal/lexer"
)

type Parser struct {
	lexer        *lexer.Lexer
	currentToken lexer.Token
}

func NewParser(l *lexer.Lexer) *Parser {
	p := &Parser{
		lexer: l,
	}
	p.advance()
	return p
}

func (p *Parser) advance() {
	p.currentToken = p.lexer.NextToken()
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
		num, err := strconv.ParseFloat(p.currentToken.Value, 64)
		if err != nil {
			return nil, &ParseError{
				Message: "Invalid number: " + p.currentToken.Value,
				Pos:     p.lexer.Pos,
			}
		}
		p.advance()
		return JSONNumber(num), nil
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

func (p *Parser) parseObject() (JSONObject, error) {
	obj := JSONObject{}

	if err := p.expectToken(lexer.LBrace); err != nil {
		return nil, err
	}

	if p.currentToken.Type == lexer.RBrace {
		p.advance()
		return obj, nil
	}

	for {
		if p.currentToken.Type != lexer.String {
			return nil, &ParseError{
				Message: "Expected string key in object",
				Pos:     p.lexer.Pos,
			}
		}
		key := p.currentToken.Value
		p.advance()

		if err := p.expectToken(lexer.Colon); err != nil {
			return nil, err
		}

		value, err := p.parseValue()
		if err != nil {
			return nil, err
		}
		obj[key] = value

		if p.currentToken.Type == lexer.RBrace {
			p.advance()
			break
		}
		if p.currentToken.Type != lexer.Comma {
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
	arr := JSONArray{}

	if err := p.expectToken(lexer.LBracket); err != nil {
		return nil, err
	}

	if p.currentToken.Type == lexer.RBracket {
		p.advance()
		return arr, nil
	}

	for {
		value, err := p.parseValue()
		if err != nil {
			return nil, err
		}
		arr = append(arr, value)

		if p.currentToken.Type == lexer.RBracket {
			p.advance()
			break
		}
		if p.currentToken.Type != lexer.Comma {
			return nil, &ParseError{
				Message: "Expected ',' or ']' in array",
				Pos:     p.lexer.Pos,
			}
		}
		p.advance()
	}

	return arr, nil
}
