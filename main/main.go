package main

import (
	"fmt"
	. "interpreter2/parser"
)

func main() {
	code := `
	let a = null;
	let b = null;
	`
	parser := NewParser(code)
	ast := parser.ParseProgram()
	fmt.Println(FormatStruct(ast))
}
