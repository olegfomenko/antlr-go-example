package main

import (
	"antrl-go-example/internal"
	expr "antrl-go-example/internal/antlr"
	"github.com/antlr4-go/antlr/v4"
	"os"
)

func main() {
	contents, err := os.ReadFile("./text.expr")

	visitor := internal.NewCalcVisitor()

	err = parseInput(visitor, string(contents))
	if err != nil {
		panic(err)
	}
}

func parseInput(visitor expr.ExprVisitor, input string) error {
	inputStream := antlr.NewInputStream(input)
	lexer := expr.NewExprLexer(inputStream)
	tokens := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	p := expr.NewExprParser(tokens)
	tree := p.Prog()
	visitor.Visit(tree)
	return nil
}
