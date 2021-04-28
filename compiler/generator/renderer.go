package generator

import "strings"

func (self *Generator) CreateRenderer(fragment Fragment) string {
	template := `const $fragmentName$ = (target, context$anchor$) => {
		$initStatements$
		return {
			$context$
			update: ($contextChain$) => {
				$updateStatements$
			},
			teardown: () => {
				$teardownStatements$
			}
		}
	}`
	anchor := ""
	if fragment.UseAnchor {
		anchor = `, anchor`
	}
	context := ""
	if fragment.HasContext {
		context = `setContext: (context) => {
			currentContext = context
		},
		getContext: () => {
			return currentContext
		},`
	}
	variables := map[string]string{
		"fragmentName": fragment.Name,
		"anchor": anchor,
		"initStatements": strings.Join(fragment.InitStatements[:], "\n\n"),
		"context": context,
		"contextChain": strings.Join(fragment.ContextChain, ", "),
		"updateStatements": strings.Join(fragment.UpdateStatments, "\n\n"),
		"teardownStatements": strings.Join(fragment.TeardownStatements, "\n\n"),
	}
	renderer := self.BuildString(template, variables)

	return renderer
}