package generator

import (
	"elljo/compiler/js-parser/ast"
	"elljo/compiler/parser"
	"strconv"
	"strings"
)

func (self *Generator) VisitLoop(children parser.Entry, current *Fragment) *Fragment {
	self.eachCounter++
	name := "loop_" + strconv.Itoa(self.eachCounter)
	renderer := "renderLoop_" + strconv.Itoa(self.eachCounter)

	template := `var $name$_anchor = document.createComment('#each $expression$');
				$target$.appendChild($name$_anchor);
				var $name$_iterations = [];
				const $name$_fragment = document.createDocumentFragment();
				this.contexts['$name$'] = [];`
	variables := map[string]string{
		"name":       name,
		"expression": children.ExpressionSource,
		"target":     current.Target,
	}
	createStatement := Statement{
		source:   self.BuildString(template, variables),
		mappings: [][]int{{}, {}, {}, {}, {}},
	}

	current.InitStatements = append(current.InitStatements, createStatement)

	teardownStatementTemplate := `for(let i = 0; i < $name$_iterations.length; i++) {
		$name$_iterations[i].teardown();
	}`

	teardownStatement := Statement{
		source:   self.BuildString(teardownStatementTemplate, variables),
		mappings: [][]int{{}, {}, {}},
	}

	current.TeardownStatements = append(current.TeardownStatements, teardownStatement)

	for _, declaration := range children.Expression.Body {
		if id, ok := declaration.(*ast.ExpressionStatement); ok && id != nil {
			variableName := children.ExpressionSource[id.Index0():id.Index1()]

			updateStatementTemplate := `if(oldState && oldState.$variableName$ && oldState.$variableName$.length > context.$variableName$.length) {
				var arrayDiff = this.utils.diffArray(context.$variableName$.length > oldState.$variableName$.length ? 
				context.$variableName$ : oldState.$variableName$, context.$variableName$.length > oldState.$variableName$.length ?
				oldState.$variableName$ : context.$variableName$);
				for(var i = 0; i < arrayDiff.length; i++) {
					$name$_iterations[arrayDiff[i].index].teardown();
				}
				for(var i = 0; i < arrayDiff.length; i++) {
					$name$_iterations.splice(arrayDiff[i].index, 1)
				}
			}
			for(var i = 0; i < context.$variableName$.length; i++) {
				if(!$name$_iterations[i]) {
					this.contexts['$name$'][i] = { $context$: context.$variableName$[i] };
					$name$_iterations[i] = $renderer$($name$_fragment, this.contexts['$name$'][i]);
					$name$_iterations[i].update($contextChain$, context.$variableName$[i], this.contexts['$name$'][i]);
				}
				const iteration = $name$_iterations[i];
				if(iteration.getContext() !== this.contexts['$name$'][i]) {
					this.contexts['$name$'][i] = { $context$: context.$variableName$[i] };
					iteration.setContext(this.contexts['$name$'][i])
				}
				iteration.update($contextChain$, context.$variableName$[i], this.contexts['$name$'][i]);
			}
			$name$_anchor.parentNode.insertBefore($name$_fragment, $name$_anchor);
			$name$_iterations.length = context.$variableName$.length;`

			variables := map[string]string{
				"variableName": variableName,
				"name":         name,
				"context":      children.Context,
				"renderer":     renderer,
				"contextChain": strings.Join(current.ContextChain, ", "),
			}

			updateStatement := Statement{
				source: self.BuildString(updateStatementTemplate, variables),
				mappings: [][]int{{}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {},
					{}, {}, {}, {}, {}, {}, {}, {}, {}, {}},
			}

			current.UpdateStatments = append(current.UpdateStatments, updateStatement)
		}
	}
	current.ContextChain = append(current.ContextChain, children.Context)
	return &Fragment{
		UseAnchor:          true,
		Name:               renderer,
		Target:             "target",
		ContextChain:       current.ContextChain,
		InitStatements:     []Statement{{source: "var currentContext = context;", mappings: [][]int{{}}}},
		UpdateStatments:    []Statement{},
		TeardownStatements: []Statement{},
		Counters: FragmentCounter{
			Text:    0,
			Anchor:  0,
			Element: 0,
		},
		Parent:     current,
		HasContext: true,
	}
}

func (self *Generator) VisitLoopAfter(current *Fragment) {
	self.renderers = append(self.renderers, self.CreateRenderer(*current))
	current = current.Parent
}
