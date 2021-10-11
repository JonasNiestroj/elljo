package ssr

import (
	"elljo/compiler/parser"
)

func (self *SSR) VisitText(parser parser.Parser, children parser.Entry) {
	self.Render += children.Data
}
