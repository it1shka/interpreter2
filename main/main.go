package main

import (
	. "interpreter2/evaluator"
	. "interpreter2/parser"
)

func main() {
	code := `
	let n = 1000;
	let sieve = ([false] * 2) + ([true] * (n - 2));
	
	let i = 2;
	for i < n {
		if sieve[i] {
			echo(i);
			let j = i * 2;
			for j < n {
				sieve[j] = false;
				j = j + i;
			}
		}
		i = i + 1;
	}
	`
	parser := NewParser(code)
	ast := parser.ParseProgram()
	// fmt.Println(FormatStruct(ast))
	eval := NewEvaluator()
	eval.Eval(ast)
}
