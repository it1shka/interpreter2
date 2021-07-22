package main

import (
	. "interpreter2/evaluator"
	. "interpreter2/parser"
)

func main() {
	code := `
	let greet = fn {
		echo("Hello world!");	
	};
	let result = greet();
	echo(result);
	`
	parser := NewParser(code)
	ast := parser.ParseProgram()
	// fmt.Println(FormatStruct(ast))
	eval := NewEvaluator()
	eval.Eval(ast)
}
