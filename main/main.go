package main

import (
	. "interpreter2/evaluator"
	. "interpreter2/parser"
	"os"
)

func RunCode(code string) {
	parser := NewParser(code)
	ast := parser.ParseProgram()
	eval := NewEvaluator()
	eval.Eval(ast)
}

func RunFromFile() {
	args := os.Args
	path := args[0]
	code := LoadFile(path)
	RunCode(code)
}

func main() {
	RunFromFile()
}
