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

func size(args []OBJECT) OBJECT {
	arg := args[0]
	if arg.getType() == ARRAY_TYPE {
		arr := arg.(*Array)
		return Int(len(*arr))
	}
	if arg.getType() == STRING_TYPE {
		return Int(len(arg.(String)))
	}
	return Int(0)
}

func push(args []OBJECT) OBJECT {
	if args[0].getType() != ARRAY_TYPE {
		return Null_
	}
	arr, elem := args[0].(*Array), args[1]
	*arr = append(*arr, elem)
	return elem
}

func pop(args []OBJECT) OBJECT {
	if args[0].getType() != ARRAY_TYPE {
		return Null_
	}
	arr := args[0].(*Array)
	if len(*arr) <= 0 {
		return Null_
	}
	elem := (*arr)[len(*arr)-1]
	*arr = (*arr)[:len(*arr)-1]
	return elem
}
