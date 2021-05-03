package generator

import "strings"

type Renderer struct {
	source   string
	mappings [][]int
}

func (self *Generator) CreateRenderer(fragment Fragment) Renderer {
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
	mappings := [][]int{{}}

	var initStatements []string

	for _, initStatement := range fragment.InitStatements {
		mappings = append(mappings, initStatement.mappings...)
		initStatements = append(initStatements, initStatement.source)
	}

	mappings = append(mappings, []int{})

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
		mappings = append(mappings, []int{}, []int{}, []int{}, []int{}, []int{}, []int{})
	}

	mappings = append(mappings, []int{})

	var updateStatements []string

	for _, updateStatement := range fragment.UpdateStatments {
		mappings = append(mappings, updateStatement.mappings...)
		updateStatements = append(updateStatements, updateStatement.source)
	}

	mappings = append(mappings, []int{}, []int{})

	var teardownStatements []string

	for _, teardownStatement := range fragment.TeardownStatements {
		mappings = append(mappings, teardownStatement.mappings...)
		teardownStatements = append(teardownStatements, teardownStatement.source)
	}

	mappings = append(mappings, []int{}, []int{}, []int{})

	variables := map[string]string{
		"fragmentName":       fragment.Name,
		"anchor":             anchor,
		"initStatements":     strings.Join(initStatements, "\n\n"),
		"context":            context,
		"contextChain":       strings.Join(fragment.ContextChain, ", "),
		"updateStatements":   strings.Join(updateStatements, "\n\n"),
		"teardownStatements": strings.Join(teardownStatements, "\n\n"),
	}
	renderer := self.BuildString(template, variables)

	return Renderer{
		source:   renderer,
		mappings: mappings,
	}
}
