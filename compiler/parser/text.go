package parser

import "strings"

func Text(parser *Parser) {
	start := parser.Index

	data := ""

	for parser.Index < len(parser.Template) && !parser.Matches("<") && !parser.Matches("{{") {
		if []byte(string(parser.Template[parser.Index]))[0] == 10 {
			data += " "
			parser.Index++
			for parser.Index < len(parser.Template) {
				if []byte(string(parser.Template[parser.Index]))[0] == 32 {
					parser.Index++
				} else {
					break
				}
			}
			continue
		}
		data += string(parser.Template[parser.Index])
		parser.Index++
	}

	if len(data) == 0 || data == "\n" || []byte(data)[0] == 10 {

		return
	}

	data = strings.ReplaceAll(data, "\n", "")

	parser.Entries[len(parser.Entries)-1].Children = append(parser.Entries[len(parser.Entries)-1].Children, &Entry{
		StartIndex: start,
		EndIndex:   parser.Index,
		EntryType:  "Text",
		Data:       data,
		Line:       parser.currentLine,
	})
}
