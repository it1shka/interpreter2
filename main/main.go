package main

import (
	. "interpreter2/evaluator"
	. "interpreter2/parser"
)

func main() {
	code := `
	let arr = [1, 2, 3];
	pop(arr);
	echo(arr);
	`
	parser := NewParser(code)
	ast := parser.ParseProgram()
	// fmt.Println(FormatStruct(ast))
	eval := NewEvaluator()
	eval.Eval(ast)
}
