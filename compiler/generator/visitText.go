package generator

import "elljo/compiler/parser"

func (self *Generator) VisitText(children parser.Entry, current *Fragment) *Fragment {
	createStatement := current.Target + ".appendChild(document.createTextNode('" + children.Data + "'))"

	current.InitStatements = append(current.InitStatements, createStatement)
	return current
}
