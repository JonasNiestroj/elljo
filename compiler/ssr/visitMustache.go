package ssr

import (
	"elljo/compiler/parser"
)

func (self *SSR) VisitMustache(parser parser.Parser, children parser.Entry) {
	self.Render += "${Main.data." + children.Parameter + "}"
}
