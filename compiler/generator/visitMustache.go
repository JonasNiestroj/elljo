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
			contextVariable := variableName
			if len(current.ContextChain) == 1 {
				contextVariable = "context." + contextVariable
			}
			updateStatementTemplate := `if($variable$ !== $name$_value) {
				$name$_value = $variable$;
				$name$.data = $name$_value;
			}`

			variables := map[string]string{
				"name":     name,
				"variable": contextVariable,
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
