package parser

import (
	"elljo/compiler/utils"
	"regexp"
)

func ReadStyle(parserInstance *Parser, start int, isGlobal bool) StyleSource {
	styleStart := parserInstance.Index
	pattern, _ := regexp.Compile("</style>")
	parserInstance.ReadUntil(pattern)

	source := parserInstance.Template[styleStart:parserInstance.Index]
	cssResult := Result{}

	// Only parse the stylesheet if we are in a scoped context
	if !isGlobal {
		cssResult = ParseStyleSheet(source)
	}

	parserInstance.Index += 8
	return StyleSource{
		StartIndex: start,
		EndIndex:   parserInstance.Index,
		Rules:      cssResult.Rules,
		Source:     source,
		Id:         utils.RandString(8),
	}
}
