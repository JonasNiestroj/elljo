package generator

import "strings"

type Renderer struct {
	source   string
	mappings [][]int
}

func (self *Generator) CreateRenderer(fragment Fragment) Renderer {
	template := `const $fragmentName$ = (target, anchor $updateContextChainParam$) => {
		$initStatements$
		return {
			update: ($updateContextChain$) => {
				$updateStatements$
			},
			teardown: () => {
				$teardownStatements$
			}
		}
	}`
	mappings := [][]int{{}}

	updateContextChainParam := fragment.UpdateContextChain
	if updateContextChainParam != "" {
		updateContextChainParam = ", " + updateContextChainParam
	}

	var initStatements []string

	for _, initStatement := range fragment.InitStatements {
		mappings = append(mappings, initStatement.mappings...)
		initStatements = append(initStatements, initStatement.source)
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
		"fragmentName":            fragment.Name,
		"initStatements":          strings.Join(initStatements, "\n\n"),
		"contextChain":            strings.Join(fragment.ContextChain, ", "),
		"updateStatements":        strings.Join(updateStatements, "\n\n"),
		"teardownStatements":      strings.Join(teardownStatements, "\n\n"),
		"updateContextChain":      fragment.UpdateContextChain,
		"updateContextChainParam": updateContextChainParam,
	}
	renderer := self.BuildString(template, variables)

	return Renderer{
		source:   renderer,
		mappings: mappings,
	}
}
