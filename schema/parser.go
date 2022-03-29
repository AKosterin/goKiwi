package schema

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var tokenRegex = regexp.MustCompile(`((?:-|\b)\d+\b|[=;{}]|\[\]|\[deprecated\]|\b[A-Za-z_][A-Za-z0-9_]*\b|\/\/.*|\s+)`)
var whitespace = regexp.MustCompile(`^\/\/.*|\s+$`)
var identifier = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)
var integer = regexp.MustCompile(`^-?\d+$`)
var semicolon = regexp.MustCompile(`^;$`)
var leftBrace = regexp.MustCompile(`^\{$`)
var rightBrace = regexp.MustCompile(`^\}$`)
var arrayToken = regexp.MustCompile(`^\[\]$`)
var equals = regexp.MustCompile(`^=$`)
var endOfFile = regexp.MustCompile(`^$`)
var deprecatedToken = regexp.MustCompile(`^\[deprecated\]$`)

type token struct {
	text   string
	line   int
	column int
}

func tokenize(text string) ([]*token, error) {
	parts := tokenRegex.FindAllString(text, -1)
	tokens := make([]*token, 0)
	line := 0
	column := 0

	for _, part := range parts {

		if !whitespace.MatchString(part) {
			tokens = append(tokens, &token{
				text:   part,
				line:   line + 1,
				column: column + 1,
			})
		}

		lines := strings.Split(part, "\n")
		if len(lines) > 1 {
			column = 0
		}
		line += len(lines) - 1
		column += len(lines[len(lines)-1])
	}

	tokens = append(tokens, &token{
		text:   "",
		line:   line,
		column: column,
	})

	return tokens, nil
}

type parser struct {
	tokens []*token
	index  int
}

func (p *parser) current() *token {
	return p.tokens[p.index]
}

func (p *parser) eat(reg *regexp.Regexp) bool {
	if !reg.MatchString(p.current().text) {
		return false
	}
	p.index++
	return true
}

func (p *parser) expect(reg *regexp.Regexp, expected string) error {
	if !p.eat(reg) {
		token := p.current()
		return fmt.Errorf("expected %s but found %q (%d:%d)", expected, token.text, token.line, token.column)
	}
	return nil
}

func (p *parser) unexpectedToken() error {
	token := p.current()
	return fmt.Errorf("unexpected token %q (%d:%d)", token.text, token.line, token.column)
}

func (p *parser) parse() (*Schema, error) {
	schema := &Schema{
		PackageName: "",
		Definitions: make([]*Definition, 0),
	}

	if p.current().text == "package" {
		p.index++
		schema.PackageName = p.current().text
		if err := p.expect(identifier, "identifier"); err != nil {
			return schema, err
		}
		if err := p.expect(semicolon, "\";\""); err != nil {
			return schema, err
		}
	}

	for p.index < len(p.tokens) && !p.eat(endOfFile) {
		definition := &Definition{
			Fields: make([]*Field, 0),
		}
		switch p.current().text {
		case "enum":
			definition.Kind = KIND_ENUM
		case "message":
			definition.Kind = KIND_MESSAGE
		case "struct":
			definition.Kind = KIND_STRUCT
		default:
			return schema, p.unexpectedToken()
		}
		p.index++

		definition.Name = p.current().text
		if err := p.expect(identifier, "identifier"); err != nil {
			return schema, err
		}
		if err := p.expect(leftBrace, "\"{\""); err != nil {
			return schema, err
		}
		for !p.eat(rightBrace) {
			field := &Field{}
			if definition.Kind != KIND_ENUM {
				field.FieldType = p.current().text

				if err := p.expect(identifier, "identifier"); err != nil {
					return schema, err
				}

				field.IsArray = p.eat(arrayToken)
			}

			field.line = p.current().line
			field.column = p.current().column
			field.Name = p.current().text
			if err := p.expect(identifier, "identifier"); err != nil {
				return schema, err
			}

			if definition.Kind != KIND_STRUCT {
				if err := p.expect(equals, "\"=\""); err != nil {
					return schema, err
				}
				v, err := strconv.ParseUint(p.current().text, 10, 32)
				if err != nil {
					token := p.current()
					return schema, fmt.Errorf("invalid integer %q (%d:%d)", token.text, token.line, token.column)
				}
				field.Value = new(uint32)
				*field.Value = uint32(v)
				if err := p.expect(integer, "integer"); err != nil {
					return schema, err
				}
			}

			depTok := p.current()
			if p.eat(deprecatedToken) {
				if definition.Kind != KIND_MESSAGE {
					return schema, fmt.Errorf("cannot deprecate this field (%d:%d)", depTok.line, depTok.column)
				}
				field.IsDepricated = true
			}

			if err := p.expect(semicolon, "\";\""); err != nil {
				return schema, err
			}
			definition.Fields = append(definition.Fields, field)
		}
		schema.Definitions = append(schema.Definitions, definition)
	}
	return schema, nil
}

func Parse(text string) (*Schema, error) {
	tokens, err := tokenize(text)
	if err != nil {
		return nil, err
	}
	return (&parser{
		tokens: tokens,
	}).parse()
}
