package generator

import (
	"elljo/compiler/parser"
	"elljo/compiler/utils"
	"strconv"
)

func (self *Generator) VisitMustache(parserInstance parser.Parser, children parser.Entry, current *Fragment) *Fragment {
	self.textCounter++
	name := "text_" + strconv.Itoa(self.textCounter)

	createStatementTemplate := "var $name$_value = '';"

	if len(children.LoopIndices) == 0 {
		createStatementTemplate += `
			var $name$ = document.createTextNode('');
			$target$.appendChild($name$);
		`
	} else {
		createStatementTemplate += `
				var $name$ = $initHtml`
		for _, index := range children.LoopIndices {
			createStatementTemplate += ".childNodes[" + strconv.Itoa(index) + "]"
		}
	}

	variables := map[string]string{
		"name":   name,
		"target": current.Target,
	}
	createStatement := Statement{
		source:   utils.BuildString(createStatementTemplate, variables),
		mappings: [][]int{{}, {}, {}},
	}

	current.InitStatements = append(current.InitStatements, createStatement)

	variable := children.ExpressionSource

	var updateStatementTemplate string
	if current.UpdateContextChain != "" {
		updateStatementTemplate += `if(!$name$_value || $variable$ !== $name$_value) {`
	} else {
		updateStatementTemplate += `if((this.$variable$IsDirty || !$name$_value) && $variable$ !== $name$_value) {`
	}

	if len(children.LoopIndices) > 0 {
		updateStatementTemplate += `
					$name$_value = $variable$;
					$name$.textContent = $name$_value;
				}`
	} else {
		updateStatementTemplate += `
					$name$_value = $variable$;
					$name$.data = $name$_value;
				}`
	}

	variables = map[string]string{
		"name":     name,
		"variable": variable,
	}

	updateStatement := Statement{
		source:   utils.BuildString(updateStatementTemplate, variables),
		mappings: [][]int{{}, {0, 0, children.Line, 0}, {}, {}},
	}
	current.UpdateStatments = append(current.UpdateStatments, updateStatement)

	return current
}
