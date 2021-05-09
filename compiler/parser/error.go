package parser

type Error struct {
	Line        int    `json:"line"`
	Message     string `json:"message"`
	StartColumn int    `json:"startColumn"`
	EndColumn   int    `json:"endColumn"`
}

func (self *Parser) Error(errorMessage string) {
	self.Errors = append(self.Errors, Error{
		Message:     errorMessage,
		Line:        self.currentLine + 1,
		StartColumn: self.PossibleErrorIndex - self.lineStartIndex,
		EndColumn:   self.Index - self.lineStartIndex,
	})
}
