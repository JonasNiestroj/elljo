package parser

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

type CSSParser struct {
	Char     rune
	Index    int
	Length   int
	Template string
}

type Rule struct {
	StartIndex int
	EndIndex   int
	Selector   string
}

type Result struct {
	Rules []Rule
}

func ParseStyleSheet(source string) Result {
	parser := CSSParser{Template: source, Length: len(source)}
	return parser.ParseStyleSheet()
}

func (self *CSSParser) ParseStyleSheet() Result {
	result := Result{}
	start := self.Index
	for {
		if self.Index == self.Length {
			break
		}

		self.read()
		self.skipWhiteSpace()
		if self.Char == '(' {
			self.readUntil(")")
		} else if self.Char == '[' {
			self.readUntil("]")
		}

		if self.Char == ',' || self.Char == '{' {
			selector := self.Template[start : self.Index-1]
			selector = strings.TrimRight(selector, " ")
			result.Rules = append(result.Rules, Rule{
				StartIndex: start,
				EndIndex:   self.Index - 1,
				Selector:   selector,
			})

			if self.Char == '{' {
				self.readUntil("}")
			}
			start = self.Index
		}
	}
	return result
}

func (self *CSSParser) peek() rune {
	if self.Index+1 < self.Length {
		return rune(self.Template[self.Index+1])
	}
	return -1
}

func (self *CSSParser) skipWhiteSpace() {
	for {
		switch self.Char {
		case ' ', '\t', '\f', '\v', '\u00a0', '\ufeff':
			self.read()
			continue
		case '\r':
			if self.peek() == '\n' {
				self.read()
			}
			fallthrough
		case '\u2028', '\u2029', '\n':
			self.read()
			continue
		}
		if self.Char >= utf8.RuneSelf {
			if unicode.IsSpace(self.Char) {
				self.read()
				continue
			}
		}
		break
	}
}

func (self *CSSParser) read() {
	if self.Index < self.Length {
		chr, width := rune(self.Template[self.Index]), 1
		self.Index += width
		self.Char = chr
	} else {
		self.Index = self.Length
		self.Char = -1
	}
}

func (self *CSSParser) readUntil(until string) string {
	var match = strings.Index(self.Template[self.Index:len(self.Template)], until)

	if match == -1 {
		str := self.Template[self.Index:len(self.Template)]
		return str
	}

	start := self.Index
	self.Index += match

	str := self.Template[start:self.Index]
	return str
}
