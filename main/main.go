package main

import (
	"fmt"
	. "interpreter2/parser"
)

func main() {
	code := `
		a = b = c = 5;
	`
	parser := NewParser(code)
	ast := parser.ParseProgram()
	fmt.Println(FormatStruct(ast))
}
