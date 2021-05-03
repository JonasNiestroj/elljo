package generator

import (
	"elljo/compiler/js-parser/ast"
	"elljo/compiler/parser"
	"strconv"
)

func (self *Generator) VisitIf(children parser.Entry, current *Fragment) *Fragment {
	self.ifCounter++
	name := "ifBlock_" + strconv.Itoa(self.ifCounter)
	renderer := "renderIfBlock_" + strconv.Itoa(self.ifCounter)

	template := `var $name$_anchor = document.createComment('#if $expression$');
				$target$.appendChild($name$_anchor);
				var $name$ = null`

	variables := map[string]string{
		"name":       name,
		"expression": children.ExpressionSource,
		"target":     current.Target,
	}

	createStatement := Statement{
		source:   self.BuildString(template, variables),
		mappings: [][]int{{}, {}, {}},
	}

	current.InitStatements = append(current.InitStatements, createStatement)

	teardownStatementSource := "if(" + name + ") " + name + ".teardown();"

	teardownStatement := Statement{
		source:   teardownStatementSource,
		mappings: [][]int{{}},
	}

	current.TeardownStatements = append(current.TeardownStatements, teardownStatement)
	for _, declaration := range children.Expression.Body {
		if id, ok := declaration.(*ast.ExpressionStatement); ok && id != nil {
			variableName := children.ExpressionSource[id.Index0():id.Index1()]
			updateStatementTemplate := `if(context.$variableName$ && !$name$){
				$name$ = $renderer$($target$, $name$_anchor);
			} else if(!context.$variableName$ && !$name$){
				$name$.teardown();
				$name$ = null;
			}
			if($name$) {
				$name$.update(context, dirtyInState, oldState);
			}`

			variables := map[string]string{
				"variableName": variableName,
				"name":         name,
				"renderer":     renderer,
				"target":       current.Target,
			}

			updateStatement := Statement{
				source:   self.BuildString(updateStatementTemplate, variables),
				mappings: [][]int{{0, 0, children.Line, 0}, {}, {}, {}, {}, {}, {}, {}, {}},
			}

			current.UpdateStatments = append(current.UpdateStatments, updateStatement)
		}
	}

	return &Fragment{
		UseAnchor:          true,
		Name:               renderer,
		Target:             "target",
		ContextChain:       current.ContextChain,
		InitStatements:     []Statement{},
		UpdateStatments:    []Statement{},
		TeardownStatements: []Statement{},
		Counters: FragmentCounter{
			Text:    0,
			Anchor:  0,
			Element: 0,
		},
		Parent: current,
	}
}

func (self *Generator) VisitIfAfter(current *Fragment) {
	self.renderers = append(self.renderers, self.CreateRenderer(*current))
	current = current.Parent
}
