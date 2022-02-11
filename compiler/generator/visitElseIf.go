package generator

import (
	"elljo/compiler/parser"
	"strconv"
)

func (self *Generator) VisitElseIf(children parser.Entry, current *Fragment) *Fragment {
	renderer := "renderElseIfBlock_" + strconv.Itoa(self.elseIfCounter)

	self.elseIfCounter++

	return &Fragment{
		UseAnchor:          true,
		Name:               renderer,
		Target:             "target",
		ContextChain:       current.ContextChain,
		InitStatements:     []Statement{},
		UpdateStatments:    []Statement{},
		TeardownStatements: []Statement{},
		Counters: &FragmentCounter{
			Text:    0,
			Anchor:  0,
			Element: 0,
		},
		Parent:             current,
		UpdateContextChain: current.UpdateContextChain,
	}
}

func (self *Generator) VisitElseIfAfter(current *Fragment) {
	self.renderers = append(self.renderers, self.CreateRenderer(*current))
	current = current.Parent
}
