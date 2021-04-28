package parser

import (
	"fmt"
	"elljo/compiler/js-parser/token"
	"sort"
)

const (
	err_UnexpectedToken = "Unexpected token %v"
	err_UnexpectedEndOfInput = "Unexpected end of input"
)

type Error struct {
	Message string
}

func (self Error) Error() string {
	return fmt.Sprintf("%s", self.Message)
}

func (self *Parser) Error(msg string, msgValues ...interface{}) *Error {
	msg = fmt.Sprintf(msg, msgValues...)
	self.Errors.Add(msg)
	return self.Errors[len(self.Errors) - 1]
}

func (self *Parser) ErrorUnexpected(chr rune) error {
	if chr == -1 {
		return self.Error(err_UnexpectedEndOfInput)
	}
	return self.Error(err_UnexpectedToken, token.ILLEGAL)
}

func (self *Parser) ErrorUnexpectedToken(tkn token.Token) error {
	switch tkn {
	case token.EOF:
		return self.Error(err_UnexpectedEndOfInput)
	}
	value := tkn.ToString()
	switch tkn {
	case token.BOOLEAN, token.NULL:
		value = self.Literal
	case token.IDENTIFIER:
		return self.Error("Unexpected identifer")
	case token.KEYWORD:
		return self.Error("Unsepected reserved word")
	case token.NUMBER:
		return self.Error("Unexpected number")
	case token.STRING:
		return self.Error("Unexpected string")
	}
	return self.Error(err_UnexpectedToken, value)
}

type ErrorList []*Error

func (self *ErrorList) Add(msg string) {
	*self = append(*self, &Error{msg})
}

func (self *ErrorList) Reset()       { *self = (*self)[0:0] }
func (self ErrorList) Len() int      { return len(self) }
func (self ErrorList) Swap(i, j int) {self[i], self[j] = self[j], self[i]}
func (self ErrorList) Less(i, j int) bool {
	return false
}

func (self ErrorList) Sort() {
	sort.Sort(self)
}

func (self ErrorList) Error() string {
	switch len(self) {
	case 0:
		return "no errors"
	case 1:
		return self[0].Error()
	}
	return fmt.Sprintf("%s (and %d more errors)", self[0].Error(), len(self) - 1)
}

func (self ErrorList) Err() error {
	if len(self) == 0 {
		return nil
	}
	return self
}