package generator

import (
	"elljo/compiler/js-parser/ast"
	"elljo/compiler/parser"
	"strconv"
)

func (self *Generator) VisitMustache(parserInstance parser.Parser, children parser.Entry, current *Fragment) *Fragment {
	self.textCounter++
	name := "text_" + strconv.Itoa(self.textCounter)

	createStatementTemplate := `var $name$ = document.createTextNode('');
								var $name$_value = '';
								$target$.appendChild($name$);`

	variables := map[string]string{
		"name":   name,
		"target": current.Target,
	}
	createStatement := Statement{
		source:   self.BuildString(createStatementTemplate, variables),
		mappings: [][]int{{}, {}, {}},
	}

	current.InitStatements = append(current.InitStatements, createStatement)

	for _, declaration := range children.Expression.Body {
		if id, ok := declaration.(*ast.ExpressionStatement); ok && id != nil {
			variableName := children.ExpressionSource[id.Index0():id.Index1()]
			variable := variableName
			var updateStatementTemplate string
			if current.UpdateContextChain != "" {
				updateStatementTemplate += `if(!$name$_value || $variable$ !== $name$_value) {`
			} else {
				updateStatementTemplate += `if((currentComponent.$variable$IsDirty || !$name$_value) && $variable$ !== $name$_value) {`
			}
			updateStatementTemplate += `
				$name$_value = $variable$;
				$name$.data = $name$_value;
			}`

			variables := map[string]string{
				"name":     name,
				"variable": variable,
			}

			updateStatement := Statement{
				source:   self.BuildString(updateStatementTemplate, variables),
				mappings: [][]int{{}, {0, 0, children.Line, 0}, {}, {}},
			}
			current.UpdateStatments = append(current.UpdateStatments, updateStatement)
		}
	}

	return current
}
