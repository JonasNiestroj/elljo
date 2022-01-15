package parser

import (
	"elljo/compiler/js-parser/token"
	"elljo/compiler/utils"
	"fmt"
	"strings"
)

const (
	err_UnexpectedToken      = "Unexpected token %v"
	err_UnexpectedEndOfInput = "Unexpected end of input"
)

func (self *Parser) Error(startIndex int, msg string, msgValues ...interface{}) {
	currentLine := strings.Count(self.Template[:self.Index], "\n") + 1
	startColumn := startIndex - strings.LastIndex(self.Template[:startIndex], "\n")
	endColumn := self.Index - strings.LastIndex(self.Template[:self.Index], "\n")

	msg = fmt.Sprintf(msg, msgValues...)
	self.Errors = append(self.Errors, utils.Error{
		Message:     msg,
		StartColumn: startColumn,
		EndColumn:   endColumn,
		Line:        currentLine,
	})
}

func (self *Parser) ErrorUnexpected(chr rune, startIndex int) {
	if chr == -1 {
		self.Error(startIndex, err_UnexpectedEndOfInput)
	}
	self.Error(startIndex, err_UnexpectedToken, token.ILLEGAL)
}

func (self *Parser) ErrorUnexpectedToken(tkn token.Token, startIndex int) {
	switch tkn {
	case token.EOF:
		self.Error(startIndex, err_UnexpectedEndOfInput)
	}
	value := tkn.ToString()
	switch tkn {
	case token.BOOLEAN, token.NULL:
		value = self.Literal
	case token.IDENTIFIER:
		self.Error(startIndex, "Unexpected identifier")
	case token.KEYWORD:
		self.Error(startIndex, "Unexpected reserved word")
	case token.NUMBER:
		self.Error(startIndex, "Unexpected number")
	case token.STRING:
		self.Error(startIndex, "Unexpected string")
	}
	self.Error(startIndex, err_UnexpectedToken, value)
}
