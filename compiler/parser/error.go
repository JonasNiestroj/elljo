package parser

type Error struct {
	Line    int    `json:"line"`
	Message string `json:"message"`
}

func (self *Parser) Error(errorMessage string) {
	self.Errors = append(self.Errors, Error{
		Message: errorMessage,
		Line:    self.currentLine,
	})
}
