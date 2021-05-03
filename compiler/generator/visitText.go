package generator

import "elljo/compiler/parser"

func (self *Generator) VisitText(children parser.Entry, current *Fragment) *Fragment {
	createStatementSource := current.Target + ".appendChild(document.createTextNode('" + children.Data + "'))"

	createStatement := Statement{
		source:   createStatementSource,
		mappings: [][]int{{}},
	}

	current.InitStatements = append(current.InitStatements, createStatement)
	return current
}
