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

	createStatement := self.BuildString(template, variables)

	current.InitStatements = append(current.InitStatements, createStatement)

	teardownStatement := "if(" + name + ") " + name + ".teardown();"

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

			current.UpdateStatments = append(current.UpdateStatments, self.BuildString(updateStatementTemplate, variables))
		}
	}

	return &Fragment{
		UseAnchor:          true,
		Name:               renderer,
		Target:             "target",
		ContextChain:       current.ContextChain,
		InitStatements:     []string{},
		UpdateStatments:    []string{},
		TeardownStatements: []string{},
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
