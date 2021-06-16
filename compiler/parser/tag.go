package parser

import (
	"elljo/compiler/js-parser/ast"
	"regexp"
	"strings"
)

type Chunk struct {
	Start      int
	End        int
	Type       string
	Data       string
	Expression *ast.Program
}

type Attribute struct {
	Name         string
	Value        string
	HasValue     bool
	IsExpression bool
	IsEvent      bool
	Expression   *ast.Program
	IsCall       bool
}

func ReadTagName(parser *Parser) string {
	var pattern, _ = regexp.Compile(`(\s|\/|>)`)
	var name = parser.ReadUntil(pattern)
	return name
}

func ReadTextBetweenQuotes(parser *Parser, quote string) string {
	expr := `(` + quote + `)`
	pattern, err := regexp.Compile(expr)
	if err != nil {
		panic(err)
	}
	text := parser.ReadUntil(pattern)
	parser.Read(quote)
	return text
}

func ReadAttributeValue(parser *Parser) string {
	isSingleQuoted := parser.Read("'")
	if isSingleQuoted {
		return ReadTextBetweenQuotes(parser, "'")
	}
	isDoubleQuoted := parser.Read("\"")
	if isDoubleQuoted {
		return ReadTextBetweenQuotes(parser, "\"")
	}
	return ""
}

func ReadAttributes(parser *Parser) []Attribute {
	pattern, _ := regexp.Compile(`(=|>|\s)`)
	var attributes []Attribute
	for {
		parser.ReadWhitespace()
		name := parser.ReadUntil(pattern)
		if name == "" {
			break
		}
		isExpression := false
		isEvent := false
		if strings.HasPrefix(name, ":") {
			isExpression = true
			name = name[1:]
		}
		if strings.HasPrefix(name, "$") {
			isEvent = true
			name = name[1:]
		}
		isValueAttribute := parser.Read("=")
		if isValueAttribute {
			value := ReadAttributeValue(parser)
			if isEvent {
				expression := ReadExpression(value)

				isCall := false
				if id, ok := expression.Body[0].(*ast.ExpressionStatement); ok && id != nil {
					if expression, ok := id.Expression.(*ast.CallExpression); ok && expression != nil {
						isCall = true
					}
				}
				attributes = append(attributes, Attribute{Name: name, HasValue: true, Value: value, IsEvent: true, Expression: expression, IsCall: isCall})
			} else if isExpression {
				expression := ReadExpression(value)

				isCall := false
				if id, ok := expression.Body[0].(*ast.ExpressionStatement); ok && id != nil {
					if expression, ok := id.Expression.(*ast.CallExpression); ok && expression != nil {
						isCall = true
					}
				}
				attributes = append(attributes, Attribute{Name: name, HasValue: true, Value: value, IsExpression: true, Expression: expression, IsCall: isCall})
			} else {
				attributes = append(attributes, Attribute{Name: name, HasValue: true, Value: `"` + value + `"`, IsEvent: isEvent})
			}
		} else {
			attributes = append(attributes, Attribute{Name: name, HasValue: false})
		}
	}

	return attributes
}

func Tag(parser *Parser) {
	parser.Index++
	var start = parser.Index

	var isClosingTag = parser.Read("/")

	var name = ReadTagName(parser)

	line := parser.currentLine

	if isClosingTag {
		parser.Read(">")

		current := parser.Entries[len(parser.Entries)-1]
		current.EndIndex = parser.Index

		length := len(parser.Entries)
		parser.Entries = parser.Entries[:length-1]

		return
	}

	attributes := ReadAttributes(parser)

	parser.ReadWhitespace()
	if name == "script" || name == "style" {
		parser.Read(">")

		if name == "script" {
			parser.ScriptSource = ReadScript(parser, parser.Index)
		}

		if name == "style" {
			isGlobal := false
			for _, attribute := range attributes {
				if attribute.Name == "global" {
					isGlobal = true
				}
			}
			parser.StyleSource = ReadStyle(parser, parser.Index, isGlobal)
		}
		return
	}

	entry := &Entry{
		StartIndex: start,
		EndIndex:   -1,
		EntryType:  "Element",
		Name:       name,
		Children:   []*Entry{},
		Attributes: attributes,
		Line:       line,
	}

	parser.Entries[len(parser.Entries)-1].Children = append(parser.Entries[len(parser.Entries)-1].Children, entry)

	isVoidElement, _ := regexp.MatchString(`^(?:area|base|br|col|command|doctype|embed|hr|img|input|keygen|link|meta|param|source|track|wbr)$`, name)

	closing := parser.Read("/") || isVoidElement

	parser.Read(">")

	if closing {
		entry.EndIndex = parser.Index
	} else {
		parser.Entries = append(parser.Entries, entry)
	}

}
