package generator

import (
	"elljo/compiler/parser"
	"elljo/compiler/utils"
	"strconv"
	"strings"
)

func WalkThroughChildrens(children *parser.Entry, initHtml string, indices []int) string {
	index := 0
	if children.Name != "" {
		children.LoopIndices = indices
		initHtml += "<" + children.Name
		for _, attribute := range children.Attributes {
			if !attribute.IsExpression && !attribute.IsEvent {
				if attribute.HasValue {
					initHtml += " " + attribute.Name + "=" + attribute.Value
				} else {
					initHtml += " " + attribute.Name
				}
			}
		}
		initHtml += ">"
	} else if children.EntryType == "MustacheTag" {
		children.LoopIndices = indices[:len(indices)-1]
	}

	for _, child := range children.Children {
		initHtml = WalkThroughChildrens(child, initHtml, append(indices, index))
		index++
	}

	if children.Name != "" {
		initHtml += "</" + children.Name + ">"
	}
	return initHtml
}

func (self *Generator) VisitLoop(children parser.Entry, current *Fragment) *Fragment {
	self.eachCounter++
	name := "loop_" + strconv.Itoa(self.eachCounter)
	renderer := "renderLoop_" + strconv.Itoa(self.eachCounter)

	self.loops = append(self.loops, renderer)

	initHtml := WalkThroughChildrens(&children, "", []int{})

	initTemplate := "let $initHtml = createFragment(`$initHtml$`).cloneNode(true);"

	variables := map[string]string{
		"initHtml": initHtml,
	}

	initStatement := Statement{
		source:   utils.BuildString(initTemplate, variables),
		mappings: [][]int{{}},
	}

	template := `var $name$_anchor = document.createComment('#each $expression$');
				$target$.appendChild($name$_anchor);
				var $name$_iterations = [];
				const $name$_fragment = document.createDocumentFragment();`
	variables = map[string]string{
		"name":       name,
		"expression": children.Parameter,
		"target":     current.Target,
	}
	createStatement := Statement{
		source:   utils.BuildString(template, variables),
		mappings: [][]int{{}, {}, {}, {}},
	}

	current.InitStatements = append(current.InitStatements, createStatement)

	teardownStatementTemplate := `for(let i = 0; i < $name$_iterations.length; i++) {
		$name$_iterations[i].teardown();
	}`

	teardownStatement := Statement{
		source:   utils.BuildString(teardownStatementTemplate, variables),
		mappings: [][]int{{}, {}, {}},
	}

	current.TeardownStatements = append(current.TeardownStatements, teardownStatement)

	updateStatementTemplate := `
				const oldState = this.oldState;
				if(oldState && oldState.$variableName$ && oldState.$variableName$.length > $variableName$.length) {
					var arrayDiff = this.utils.diffArray($variableName$.length > oldState.$variableName$.length ? 
					$variableName$ : oldState.$variableName$, $variableName$.length > oldState.$variableName$.length ?
					oldState.$variableName$ : $variableName$);
					for(var i = 0, length = arrayDiff.length; i < length; i++) {
						$name$_iterations[arrayDiff[i].index].teardown();
					}
					for(var i = 0, length = arrayDiff.length; i < length; i++) {
						$name$_iterations.splice(arrayDiff[i].index, 1)
					}
				}
				for(var i = 0, length = $variableName$.length; i < length; i++) {
					if(!$name$_iterations[i]) {
						var variable = $variableName$[i];
						$name$_iterations[i] = $renderer$($name$_anchor.parentNode, $name$_anchor, variable);
						$name$_iterations[i].update(variable);
					}
					const iteration = $name$_iterations[i];
					iteration.update($variableName$[i]);
				}
				$name$_iterations.length = $variableName$.length;`

	variables = map[string]string{
		"variableName": children.Parameter,
		"name":         name,
		"context":      children.Context,
		"renderer":     renderer,
		"contextChain": strings.Join(current.ContextChain, ", "),
	}

	updateStatement := Statement{
		source: utils.BuildString(updateStatementTemplate, variables),
		mappings: [][]int{{}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {},
			{}, {}, {}, {}, {}, {}},
	}

	current.UpdateStatments = append(current.UpdateStatments, updateStatement)

	current.ContextChain = append(current.ContextChain, children.Context)
	return &Fragment{
		UseAnchor:          false,
		Name:               renderer,
		Target:             "target",
		ContextChain:       current.ContextChain,
		InitStatements:     []Statement{initStatement},
		UpdateStatments:    []Statement{},
		TeardownStatements: []Statement{},
		Counters: &FragmentCounter{
			Text:    0,
			Anchor:  0,
			Element: 0,
		},
		Parent:             current,
		HasContext:         true,
		UpdateContextChain: children.Context,
	}
}

func (self *Generator) VisitLoopAfter(current *Fragment) {
	current.InitStatements = append(current.InitStatements, Statement{source: `
		target.appendChild($initHtml);`, mappings: [][]int{{}, {}}})
	self.renderers = append(self.renderers, self.CreateRenderer(*current))
	self.loops = self.loops[:len(self.loops)-1]
	current = current.Parent
}
