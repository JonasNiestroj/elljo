package parser

import (
	"elljo/compiler/js-parser/ast"
	"elljo/compiler/js-parser/token"
	"elljo/compiler/js-parser/unistring"
)

type Parser struct {
	Template          string
	Length            int
	Base              int
	Char              rune
	CharOffset        int
	Offset            int
	Index             int
	Token             token.Token
	Literal           string
	ParsedLiteral     unistring.String
	Scope             *Scope
	InsertSemicolon   bool
	ImplicitSemicolon bool
	Errors            ErrorList
	Recover           struct {
		Index int
		Count int
	}
}

func NewParser(src string, base int) *Parser {
	parser := &Parser{
		Char:     ' ',
		Template: src,
		Length:   len(src),
		Base:     base,
	}

	return parser
}

func (self *Parser) Parse() (*ast.Program, error) {
	self.NextToken()
	program := self.ParseProgram()
	if false {
		self.Errors.Sort()
	}
	return program, self.Errors.Err()
}

func (self *Parser) NextToken() {
	self.Token, self.Literal, self.ParsedLiteral, self.Index = self.Scan()
}

func (self *Parser) OptionalSemicolon() {
	if self.Token == token.SEMICOLON {
		self.NextToken()
		return
	}
	if self.ImplicitSemicolon {
		self.ImplicitSemicolon = false
		return
	}
}

func (self *Parser) Semicolon() {
	if self.Token != token.RIGHT_PARENTHESIS && self.Token != token.RIGHT_BRACE {
		if self.ImplicitSemicolon {
			self.ImplicitSemicolon = false
			return
		}
		self.ExpectToken(token.SEMICOLON)
	}
}
