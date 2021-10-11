package ssr

import (
	"elljo/compiler/parser"
)

func (self *SSR) VisitElement(parser parser.Parser, children parser.Entry) {

	template := "<" + children.Name

	if len(children.Attributes) > 0 {
		for _, attribute := range children.Attributes {
			if template[len(template)-1] != ' ' {
				template += " "
			}
			template += attribute.Name + "=" + attribute.Value
		}
	}

	if parser.StyleSource.Id != "" {
		template += " scope-" + parser.StyleSource.Id
	}

	template += ">"
	self.Render += template
}

func (self *SSR) VisitElementAfter(children parser.Entry) {
	template := "</" + children.Name + ">"

	self.Render += template
}
