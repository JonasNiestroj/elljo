package parser

import (
	"elljo/compiler/js-parser/ast"
	"regexp"
)

type Entry struct {
	StartIndex       int
	EndIndex         int
	EntryType        string
	Children         []*Entry
	Name             string
	Data             string
	Expression       *ast.Program
	ExpressionSource string
	Context          string
	Attributes       []Attribute
	Line             int
	HasElse          bool
	ElseIfs          []*Entry
	Else             *Entry
}

type ScriptSource struct {
	StartIndex int
	EndIndex   int
	Program    *ast.Program
	Variables  []string
	Imports    []string
}

type Parser struct {
	Index              int
	Template           string
	Entries            []*Entry
	ScriptSource       ScriptSource
	currentLine        int
	Errors             []Error
	PossibleErrorIndex int
	lineStartIndex     int
}

var (
	newLineRegex = regexp.MustCompile("(\\r\\n|\\r|\\n)")
)

func (self *Parser) Matches(str string) bool {
	var to = self.Index + len(str)
	if to > len(self.Template) {
		return false
	}
	return self.Template[self.Index:to] == str
}

func read(parser *Parser, str string) bool {
	if parser.Matches(str) {
		parser.Index += len(str)
		return true
	}
	return false
}

func (self *Parser) ReadRequired(str string) bool {
	read := read(self, str)
	if !read {
		self.Error("Expected " + str)
	}
	return read
}

func (self *Parser) Read(str string) bool {
	return read(self, str)
}

func (self *Parser) ReadWithWhitespaceRequired(str string) bool {
	start := self.Index
	for self.Index < len(self.Template) {
		var match, _ = regexp.MatchString(`\s`, string(self.Template[self.Index]))
		if match {
			self.Index++
		} else {
			break
		}
	}
	read := read(self, str)
	if !read {
		self.Error("Expected " + str)
	}
	readSource := self.Template[start:self.Index]
	newLines := len(newLineRegex.FindAllStringIndex(readSource, -1))
	if newLines > 0 {
		self.currentLine += newLines
		self.lineStartIndex = self.Index
	}

	return read
}

func (self *Parser) ReadWhitespace() {
	start := self.Index
	for self.Index < len(self.Template) {
		var match, _ = regexp.MatchString(`\s`, string(self.Template[self.Index]))
		if match {
			self.Index++
		} else {
			break
		}
	}

	str := self.Template[start:self.Index]
	newLines := len(newLineRegex.FindAllStringIndex(str, -1))
	if newLines > 0 {
		self.currentLine += newLines
		self.lineStartIndex = self.Index
	}
}

func (self *Parser) ReadUntil(pattern *regexp.Regexp) string {
	var match = pattern.FindStringIndex(self.Template[self.Index:len(self.Template)])

	if match == nil {
		str := self.Template[self.Index:len(self.Template)]
		newLines := len(newLineRegex.FindAllStringIndex(str, -1))
		if newLines > 0 {
			self.currentLine += newLines
			self.lineStartIndex = self.Index
		}
		return str
	}

	if match[0] == 0 {
		return ""
	}
	start := self.Index
	self.Index += match[0]

	str := self.Template[start:self.Index]
	newLines := len(newLineRegex.FindAllStringIndex(str, -1))
	if newLines > 0 {
		self.currentLine += newLines
		self.lineStartIndex = self.Index
	}
	return str
}

func (self *Parser) Parse() {
	emptyEntry := &Entry{
		StartIndex: 0,
		EndIndex:   -1,
		EntryType:  "fragment",
		Children:   []*Entry{},
		Name:       "html",
	}
	self.Entries = append(self.Entries, emptyEntry)

	for true {
		self.ReadWhitespace()

		if self.Matches("<") {
			Tag(self)
		} else if self.Matches("{{") {
			Mustache(self)
		} else {
			Text(self)
		}

		if self.Index >= len(self.Template)-1 {
			break
		}
	}
}
