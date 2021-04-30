package service

import (
	"bufio"
	"elljo/compiler/generator"
	"elljo/compiler/parser"
	"os"
	"strings"
)

func RunService() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		bytes := scanner.Bytes()
		message := strings.ReplaceAll(string(bytes), "\\n", "\n")
		if strings.HasPrefix(message, "compile") {
			var parserInstance = parser.Parser{
				Index:    0,
				Template: strings.Replace(message, "compile ", "", 1),
				Entries:  []*parser.Entry{},
			}
			parserInstance.Parse()

			generatorInstance := generator.Generator{}

			generated, sourcemap := generatorInstance.Generate(parserInstance, parserInstance.Template)

			output := generated + "$%&" + sourcemap

			os.Stdout.Write([]byte(output))

			break
		}
	}
}
