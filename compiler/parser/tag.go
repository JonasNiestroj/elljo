package parser

import (
	"elljo/compiler/utils"
	"github.com/JonasNiestroj/esbuild-internal/js_ast"
	"regexp"
	"strings"
)

type Attribute struct {
	Name         string
	Value        string
	HasValue     bool
	IsExpression bool
	IsEvent      bool
	Expression   js_ast.AST
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

func isExpressionACall(ast js_ast.AST) bool {
	for _, part := range ast.Parts {
		for _, stmt := range part.Stmts {
			if id, ok := stmt.Data.(*js_ast.SExpr); ok && id != nil {
				if _, ok := id.Value.Data.(*js_ast.ECall); ok {
					return true
				}
			}
		}
	}
	return false
}

func ReadAttributes(parser *Parser, entry *Entry) []Attribute {
	pattern, _ := regexp.Compile(`(/>|=|>|\s)`)
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
			if name == "xmlns" {
				entry.Namespace = value
			}
			if isEvent {
				expression := ReadExpression(value)

				isCall := isExpressionACall(expression)

				attributes = append(attributes, Attribute{Name: name, HasValue: true, Value: value, IsEvent: true, Expression: expression, IsCall: isCall})
			} else if isExpression {
				expression := ReadExpression(value)

				isCall := isExpressionACall(expression)

				attributes = append(attributes, Attribute{Name: name, HasValue: true, Value: value, IsExpression: true, Expression: expression, IsCall: isCall})
			} else {
				attributes = append(attributes, Attribute{Name: name, HasValue: true, Value: `"` + value + `"`, IsEvent: isEvent})
			}
		} else {
			attributes = append(attributes, Attribute{Name: name, HasValue: false})
		}
	}

	if entry.EntryType == "SlotElement" && len(attributes) > 1 {
		parser.Error("Only one name is allowed for the slot tag")
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

	entryType := "Element"

	if name == "slot" {
		entryType = "SlotElement"
	}

	entry := &Entry{
		StartIndex: start,
		EndIndex:   -1,
		EntryType:  entryType,
		Name:       name,
		Children:   []*Entry{},
		Attributes: []Attribute{},
		Line:       line,
		Namespace:  parser.Entries[len(parser.Entries)-1].Namespace,
	}

	entry.Attributes = ReadAttributes(parser, entry)

	if entryType == "SlotElement" {
		name := ""
		if len(entry.Attributes) > 0 {
			name = entry.Attributes[0].Name
		}

		for _, oldEntry := range parser.Slots {
			if oldEntry.Name == name {
				if name != "" {
					parser.Error("Duplicate slot with the name " + name + " detected")
				} else {
					parser.Error("Duplicate default slot detected")
				}

			}
		}

		parser.Slots = append(parser.Slots, Slot{Name: name})
	}

	// Lets initialize an empty ScriptSource to prevent errors with no script tag
	parser.ScriptSource = ScriptSource{
		StringReplacer: &utils.StringReplacer{},
	}

	parser.ReadWhitespace()
	if name == "script" || name == "style" {
		parser.Read(">")

		if name == "script" {
			parser.ScriptSource = ReadScript(parser, parser.Index)
		}

		if name == "style" {
			isGlobal := false
			for _, attribute := range entry.Attributes {
				if attribute.Name == "global" {
					isGlobal = true
				}
			}
			parser.StyleSource = ReadStyle(parser, parser.Index, isGlobal)
		}
		return
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
