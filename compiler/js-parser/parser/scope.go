package parser

import (
	"elljo/compiler/js-parser/ast"
	"elljo/compiler/js-parser/unistring"
)

type Scope struct {
	Outer           *Scope
	AllowIn         bool
	InIteration     bool
	InSwitch        bool
	InFunction      bool
	DeclarationList []ast.Declaration
	Labels          []unistring.String
}

func (self *Parser) OpenScope() {
	self.Scope = &Scope{
		Outer:   self.Scope,
		AllowIn: true,
	}
}

func (self *Parser) CloseScope() {
	self.Scope = self.Scope.Outer
}

func (self *Scope) Declare(declaration ast.Declaration) {
	self.DeclarationList = append(self.DeclarationList, declaration)
}

func (self *Scope) HasLabel(name unistring.String) bool {
	for _, label := range self.Labels {
		if label == name {
			return true
		}
	}
	if self.Outer != nil && !self.InFunction {
		return self.Outer.HasLabel(name)
	}
	return false
}
