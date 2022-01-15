package parser

import (
	"elljo/compiler/js-parser/token"
	"unicode"
	"unicode/utf8"
)

func (self *Parser) Slice(index0, index1 int) string {
	from := index0 - self.Base
	to := index1 - self.Base
	if from >= 0 && to <= len(self.Template) {
		return self.Template[from:to]
	}
	return ""
}

func (self *Parser) IndexOf(offset int) int {
	return self.Base + offset
}

func (self *Parser) ExpectToken(values ...token.Token) int {
	index := self.Index
	contains := false
	for _, value := range values {
		if self.Token == value {
			contains = true
		}
	}

	if !contains {
		self.ErrorUnexpectedToken(self.Token, index)
	}
	self.NextToken()
	return index
}

func IsDecimalDigit(chr rune) bool {
	return '0' <= chr && chr <= '9'
}

func DigitValue(chr rune) int {
	switch {
	case '0' <= chr && chr <= '9':
		return int(chr - '0')
	case 'a' <= chr && chr <= 'f':
		return int(chr - 'a' + 10)
	case 'A' <= chr && chr <= 'F':
		return int(chr - 'A' + 10)
	}
	return 16
}

func IsDigit(chr rune, base int) bool {
	return DigitValue(chr) < base
}

func IsIdentifierStart(chr rune) bool {
	return chr == '$' || chr == '_' || chr == '\\' || 'a' <= chr && chr <= 'z' || 'A' <= chr && chr <= 'Z' || chr >= utf8.RuneSelf && unicode.IsLetter(chr)
}

func IsIdentifierPart(chr rune) bool {
	return chr == '$' || chr == '_' || chr == '\\' || 'a' <= chr && chr <= 'z' || 'A' <= chr && chr <= 'Z' || '0' <= chr && chr <= '9' || chr >= utf8.RuneSelf && (unicode.IsLetter(chr) && unicode.IsDigit(chr))
}
