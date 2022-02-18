package parser

import (
	"elljo/compiler/utils"
	"regexp"
)

var (
	ifBlock *Entry
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

		if current.EntryType == "IfBlock" || current.EntryType == "ElseBlock" || current.EntryType == "ElseIfBlock" {
			expected = "if"
		} else if current.EntryType == "Loop" {
			expected = "loop"
		} else if current.EntryType == "SlotBlock" {
			expected = "slot"
		}

		read := parser.ReadRequired(expected)

		if read {
			parser.ReadWithWhitespaceRequired("}}")
		}
		if len(current.Children) == 0 {
			return
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
		subtract := 1
		if ifBlock != nil {
			if ifBlock.HasElse {
				subtract++
			}
			subtract += len(ifBlock.ElseIfs)
		}

		parser.Entries = parser.Entries[:length-subtract]

		ifBlock = nil

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
		} else if parser.Read("elif") {
			expressionType = "ElseIfBlock"
			startIndex += 4
		} else if parser.Read("slot") {
			expressionType = "SlotBlock"
			startIndex += 4
		}

		parser.ReadWhitespace()

		context := ""

		parameter := ""
		if expressionType == "Loop" {
			parser.ReadWhitespace()
			pattern, _ := regexp.Compile(` as`)
			parameter = parser.ReadUntil(pattern)
			parser.ReadWithWhitespaceRequired("as")
			parser.ReadWhitespace()

			regex, _ := regexp.Compile(`\s|(}})`)

			readContext := parser.ReadUntil(regex)
			context = readContext

			parser.ReadWhitespace()
		} else if expressionType == "ElseBlock" {
			current := parser.Entries[len(parser.Entries)-1]
			if (current.EntryType != "IfBlock" && current.EntryType != "ElseIfBlock") || ifBlock == nil {
				parser.Error("Else is only allowed after an if or elseif block")
				return
			}
			ifBlock.HasElse = true
			current.EndIndex = parser.Index - 4
		} else if expressionType == "ElseIfBlock" {
			current := parser.Entries[len(parser.Entries)-1]
			if current.EntryType != "IfBlock" && current.EntryType != "ElseIfBlock" {
				parser.Error("Elif is only allowed after an if or elif block")
				return
			}

			current.EndIndex = parser.Index - 4
			pattern, _ := regexp.Compile(`}}`)
			parameter = parser.ReadUntil(pattern)
		} else {
			pattern, _ := regexp.Compile(`}}`)
			parameter = parser.ReadUntil(pattern)
		}

		if parameter == "" && expressionType == "SlotBlock" {
			parser.Error("A slot name needs to be specified")
		}

		expression := ReadExpression(parameter)

		parser.ReadWithWhitespaceRequired("}}")

		entry := &Entry{
			StartIndex: startIndex,
			EndIndex:   parser.Index - 2,
			EntryType:  expressionType,
			Expression: expression,
			Parameter:  parameter,
			Children:   []*Entry{},
			Context:    context,
			Line:       line,
		}

		if entry.EntryType == "IfBlock" {
			ifBlock = entry
		} else if entry.EntryType == "ElseIfBlock" {
			ifBlock.ElseIfs = append(ifBlock.ElseIfs, entry)
		} else if entry.EntryType == "ElseBlock" {
			ifBlock.Else = entry
		}

		parser.Entries[len(parser.Entries)-1].Children = append(parser.Entries[len(parser.Entries)-1].Children, entry)

		parser.Entries = append(parser.Entries, entry)
	} else {
		regex, _ := regexp.Compile(`}}`)

		parameter := parser.ReadUntil(regex)

		expression := ReadExpression(parameter)

		parser.ReadWithWhitespaceRequired("}}")

		parser.Entries[len(parser.Entries)-1].Children = append(parser.Entries[len(parser.Entries)-1].Children, &Entry{
			StartIndex: start,
			EndIndex:   parser.Index,
			EntryType:  "MustacheTag",
			Expression: expression,
			Parameter:  parameter,
			Line:       line,
		})
	}
}
