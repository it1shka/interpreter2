package main

import (
	. "interpreter2/evaluator"
	. "interpreter2/parser"
)

func main() {
	code := `
	echo(1 == 1.0);
	`
	parser := NewParser(code)
	ast := parser.ParseProgram()
	// fmt.Println(FormatStruct(ast))
	eval := NewEvaluator()
	eval.Eval(ast)
}
