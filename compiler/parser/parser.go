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
	Attributes 		 []Attribute
}

type ScriptSource struct {
	StartIndex int
	EndIndex   int
	Program    *ast.Program
	Variables  []string
	Imports    []string
}

type Parser struct {
	Index        int
	Template     string
	Entries      []*Entry
	ScriptSource ScriptSource
}

func (self *Parser) Matches(str string) bool {
	var to = self.Index + len(str)
	if to > len(self.Template) {
		return false
	}
	return self.Template[self.Index:to] == str
}

func (self *Parser) Read(str string) bool {
	if self.Matches(str) {
		self.Index += len(str)
		return true
	}
	return false
}

func (self *Parser) ReadWhitespace() {
	for self.Index < len(self.Template) {
		var match, _ = regexp.MatchString(`\s`, string(self.Template[self.Index]))
		if match {
			self.Index++
		} else {
			break
		}
	}
}

func (self *Parser) ReadUntil(pattern *regexp.Regexp) string {
	var match = pattern.FindStringIndex(self.Template[self.Index:len(self.Template)])

	if match == nil {
		return self.Template[self.Index:len(self.Template)]
	}

	if match[0] == 0 {
		return ""
	}
	start := self.Index
	self.Index += match[0]

	return self.Template[start:self.Index]
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
