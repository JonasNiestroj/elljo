package main

import (
	"elljo/compiler/generator"
	"elljo/compiler/parser"
	"elljo/compiler/service"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	if os.Args[1] == "--service" {
		service.RunService()
	}
	inputFile := os.Args[1]
	data, err := ioutil.ReadFile(inputFile)
	if err != nil {
		panic(err)
	}

	var parserInstance = parser.Parser{
		Index:    0,
		Template: string(data),
		Entries:  []*parser.Entry{},
	}

	parserInstance.Parse()

	var generatorInstance = generator.Generator{}

	generated, sourcemap := generatorInstance.Generate(parserInstance, parserInstance.Template)

	indexFile := os.Args[2]
	index, err := ioutil.ReadFile(indexFile)
	if err != nil {
		panic(err)
	}
	output := strings.Replace(string(index), "{SCRIPT}", generated, 1)
	ioutil.WriteFile(os.Args[3], []byte(output), 0644)

	if len(os.Args) > 4 {
		ioutil.WriteFile(os.Args[4], []byte(sourcemap), 0644)
	}
}
