package parser

import "elljo/compiler/utils"

func (self *Parser) Error(errorMessage string) {
	self.Errors = append(self.Errors, utils.Error{
		Message:     errorMessage,
		Line:        self.currentLine + 1,
		StartColumn: self.PossibleErrorIndex - self.lineStartIndex,
		EndColumn:   self.Index - self.lineStartIndex,
	})
}
