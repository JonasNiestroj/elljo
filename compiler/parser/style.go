package parser

import (
	"elljo/compiler/utils"
	"regexp"
)

func ReadStyle(parserInstance *Parser, start int) StyleSource {
	styleStart := parserInstance.Index
	pattern, _ := regexp.Compile("</style>")
	parserInstance.ReadUntil(pattern)

	source := Spaces(styleStart) + parserInstance.Template[styleStart:parserInstance.Index]
	cssResult := ParseStyleSheet(source)

	parserInstance.Index += 8
	return StyleSource{
		StartIndex: start,
		EndIndex:   parserInstance.Index,
		Rules:      cssResult.Rules,
		Source:     source,
		Id:         utils.RandString(8),
	}
}
