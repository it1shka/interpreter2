package evaluator

import (
	"fmt"
)

func echo(args []OBJECT) OBJECT {
	for _, obj := range args {
		fmt.Println(obj.ToString())
	}
	return Null_
}
