package main

import (
	. "interpreter2/evaluator"
	. "interpreter2/parser"
)

func main() {
	code := `
	let a = "Hello";
	let b = "Hel";
	echo(a-b);
	`
	parser := NewParser(code)
	ast := parser.ParseProgram()
	// fmt.Println(FormatStruct(ast))
	eval := NewEvaluator()
	eval.Eval(ast)
}
