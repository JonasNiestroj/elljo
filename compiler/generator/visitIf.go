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

	elseName := ""
	elseRenderer := ""

	if children.HasElse {
		self.elseCounter++
		elseName = "elseBlock_" + strconv.Itoa(self.elseCounter)
		elseRenderer = "renderElseBlock_" + strconv.Itoa(self.elseCounter)
	}

	template := `var $name$_anchor = document.createComment('#if $expression$');
				$target$.appendChild($name$_anchor);
				var $name$ = null`

	if children.HasElse {
		template += `
			var $elseName$_anchor = document.createComment('#else');
			$target$.appendChild($elseName$_anchor);
			var $elseName$ = null;
		`
	}

	variables := map[string]string{
		"name":       name,
		"expression": children.ExpressionSource,
		"target":     current.Target,
		"elseName":   elseName,
	}

	createStatement := Statement{
		source:   self.BuildString(template, variables),
		mappings: [][]int{{}, {}, {}},
	}

	current.InitStatements = append(current.InitStatements, createStatement)

	teardownStatementSource := "if(" + name + ") " + name + ".teardown();"

	if children.HasElse {
		teardownStatementSource += `
			if(` + elseName + `) ` + elseName + `.teardown();`
	}

	teardownStatement := Statement{
		source:   teardownStatementSource,
		mappings: [][]int{{}},
	}

	current.TeardownStatements = append(current.TeardownStatements, teardownStatement)
	for _, declaration := range children.Expression.Body {
		if id, ok := declaration.(*ast.ExpressionStatement); ok && id != nil {
			variableName := children.ExpressionSource[id.Index0():id.Index1()]
			updateStatementTemplate := `if(context.$variableName$){
				if(!$name$) $name$ = $renderer$($target$, $name$_anchor);
			`

			if children.HasElse {
				updateStatementTemplate += `
						$name$.update(context, dirtyInState, oldState);
						if($elseName$) { 
							$elseName$.teardown();
							$elseName$ = null;
						}
					}
					else {
						if(!$elseName$) $elseName$ = $elseRenderer$($target$, $elseName$_anchor);
						if($name$) {
							$name$.teardown();
							$name$ = null;
						}
						$elseName$.update(context, dirtyInState, oldState);
					}`
			} else {
				updateStatementTemplate += `
				} else if(!context.$variableName$ && !$name$){
					$name$.teardown();
					$name$ = null;
				}`
			}

			variables := map[string]string{
				"variableName": variableName,
				"name":         name,
				"renderer":     renderer,
				"target":       current.Target,
				"elseName":     elseName,
				"elseRenderer": elseRenderer,
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
