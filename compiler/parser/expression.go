package parser

import (
	"github.com/JonasNiestroj/esbuild-internal/config"
	"github.com/JonasNiestroj/esbuild-internal/js_ast"
	"github.com/JonasNiestroj/esbuild-internal/js_parser"
	"github.com/JonasNiestroj/esbuild-internal/logger"
)

func ReadExpression(src string) js_ast.AST {
	log := logger.NewDeferLog(logger.DeferLogAll)

	source := logger.Source{Index: 0, Contents: src}

	astTree, _ := js_parser.Parse(log, source, js_parser.OptionsFromConfig(&config.Options{WriteToStdout: false}))

	return astTree
}
