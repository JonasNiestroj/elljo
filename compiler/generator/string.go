package generator

import "strings"

func (self *Generator) BuildString(template string, variables map[string]string) string {
	templateReplaced := template
	for variableName, variable := range variables {
		templateReplaced = strings.ReplaceAll(templateReplaced, "$" + variableName + "$", variable)
	}

	return templateReplaced
}
