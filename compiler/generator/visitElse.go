package generator

import (
	"elljo/compiler/parser"
	"strconv"
)

func (self *Generator) VisitElse(children parser.Entry, current *Fragment) *Fragment {
	renderer := "renderElseBlock_" + strconv.Itoa(self.elseCounter)
	name := "ifBlock_" + strconv.Itoa(self.elseCounter)

	teardownStatementSource := "if(" + name + ") " + name + ".teardown();"

	teardownStatement := Statement{
		source:   teardownStatementSource,
		mappings: [][]int{{}},
	}

	current.TeardownStatements = append(current.TeardownStatements, teardownStatement)

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

func (self *Generator) VisitElseAfter(current *Fragment) {
	self.renderers = append(self.renderers, self.CreateRenderer(*current))
	current = current.Parent
}
