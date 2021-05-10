package parser

import (
	"elljo/compiler/utils"
	"regexp"
)

func Mustache(parser *Parser) {
	start := parser.Index
	parser.PossibleErrorIndex = parser.Index
	parser.Index += 2
	parser.ReadWhitespace()
	line := parser.currentLine
	if parser.Read("/") {
		current := parser.Entries[len(parser.Entries)-1]

		expected := ""

		if current.EntryType == "IfBlock" || current.EntryType == "ElseBlock" {
			expected = "if"
		} else if current.EntryType == "Loop" {
			expected = "loop"
		}

		read := parser.ReadRequired(expected)
		if read {
			parser.ReadWithWhitespaceRequired("}}")
		}

		firstChild := current.Children[0]
		lastChild := current.Children[len(current.Children)-1]

		if current.StartIndex-1 >= 0 {
			charBefore := string(parser.Template[current.StartIndex-1])
			charAfter := ""
			if parser.Index < len(parser.Template) {
				charAfter = string(parser.Template[parser.Index])
			}

			if charBefore == " " {
				firstChild.Data = utils.TrimStart(firstChild.Data)
			}

			if charAfter == " " {
				lastChild.Data = utils.TrimEnd(lastChild.Data)
			}
		}

		current.EndIndex = parser.Index
		length := len(parser.Entries)
		parser.Entries = parser.Entries[:length-1]
	} else if parser.Read("#") {
		expressionType := ""
		startIndex := parser.Index
		if parser.Read("if") {
			expressionType = "IfBlock"
			startIndex += 2
		} else if parser.Read("loop") {
			expressionType = "Loop"
			startIndex += 4
		} else if parser.Read("else") {
			expressionType = "ElseBlock"
			startIndex += 4
		}

		parser.ReadWhitespace()

		context := ""

		expressionSource := ""
		if expressionType == "Loop" {
			parser.ReadWhitespace()
			pattern, _ := regexp.Compile(` as`)
			expressionSource = parser.ReadUntil(pattern)
			parser.ReadWithWhitespaceRequired("as")
			parser.ReadWhitespace()

			regex, _ := regexp.Compile(`\s|(}})`)

			readContext := parser.ReadUntil(regex)
			context = readContext

			parser.ReadWhitespace()
		} else if expressionType == "ElseBlock" {
			current := parser.Entries[len(parser.Entries)-1]
			if current.EntryType != "IfBlock" {
				parser.Error("Else is only allowed after an if block")
				return
			}
			current.HasElse = true
			current.EndIndex = parser.Index - 4
		} else {
			pattern, _ := regexp.Compile(`}}`)
			expressionSource = parser.ReadUntil(pattern)
		}

		expression := ReadExpression(expressionSource)

		parser.ReadWithWhitespaceRequired("}}")

		entry := &Entry{
			StartIndex:       startIndex,
			EndIndex:         parser.Index - 2,
			EntryType:        expressionType,
			Expression:       expression,
			ExpressionSource: expressionSource,
			Children:         []*Entry{},
			Context:          context,
			Line:             line,
		}

		parser.Entries[len(parser.Entries)-1].Children = append(parser.Entries[len(parser.Entries)-1].Children, entry)

		parser.Entries = append(parser.Entries, entry)
	} else {
		regex, _ := regexp.Compile(`}}`)

		expressionSource := parser.ReadUntil(regex)

		expression := ReadExpression(expressionSource)

		parser.ReadWithWhitespaceRequired("}}")

		parser.Entries[len(parser.Entries)-1].Children = append(parser.Entries[len(parser.Entries)-1].Children, &Entry{
			StartIndex:       start,
			EndIndex:         parser.Index,
			EntryType:        "MustacheTag",
			Expression:       expression,
			ExpressionSource: expressionSource,
			Line:             line,
		})
	}
}
