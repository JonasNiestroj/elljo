package utils

import (
	"strings"
)

func BuildString(template string, variables map[string]string) string {
	templateReplaced := template
	for variableName, variable := range variables {
		templateReplaced = strings.ReplaceAll(templateReplaced, "$"+variableName+"$", variable)
	}

	return templateReplaced
}

// CountLines returns the lines count from the given string
func CountLines(str string) int {
	lines := 0

	for _, r := range str {
		if r == '\n' {
			lines++
		}
	}

	// If the strings is not ending with a \n, we need to add an extra line
	if len(str) > 0 && !strings.HasSuffix(str, "\n") {
		lines++
	}

	return lines
}
