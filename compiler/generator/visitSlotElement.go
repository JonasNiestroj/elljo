package generator

import (
	"elljo/compiler/parser"
	"elljo/compiler/utils"
)

func (self *Generator) VisitSlotElement(children parser.Entry, current *Fragment) *Fragment {
	slotName := "default"

	if len(children.Attributes) > 0 {
		slotName = children.Attributes[0].Name
	}

	createStatementTemplate := `
		this.$slotTargets['$slotName$'] = $target$;
		if(this.$slots['$slotName$']) {
			var slotFragment = document.createDocumentFragment();
			this.$slots['$slotName$']().render(slotFragment);
			$target$.appendChild(slotFragment);
		}
	`

	variables := map[string]string{
		"target":   current.Target,
		"slotName": slotName,
	}

	createStatement := Statement{
		source:   utils.BuildString(createStatementTemplate, variables),
		mappings: [][]int{{}, {}, {}},
	}

	current.InitStatements = append(current.InitStatements, createStatement)

	return current
}
