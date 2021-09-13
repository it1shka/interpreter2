package evaluator

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func typeof(args []OBJECT) OBJECT {
	return Int(args[0].getType())
}

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

func shift(args []OBJECT) OBJECT {
	if args[0].getType() != ARRAY_TYPE || len(*args[0].(*Array)) < 1 {
		return Null_
	}
	first := (*args[0].(*Array))[0]
	*args[0].(*Array) = (*args[0].(*Array))[1:]
	return first
}

func input(args []OBJECT) OBJECT {
	if len(args) > 0 {
		fmt.Println(args[0].ToString())
	}
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	text = strings.TrimRight(text, "\r\n")
	return String(text)
}

func integer(args []OBJECT) OBJECT {
	switch tt := args[0].(type) {
	case Int:
		return tt
	case Float:
		return Int(tt)
	case *Bool:
		if tt == False_ {
			return Int(0)
		} else {
			return Int(1)
		}
	case *Null:
		return Int(0)
	case String:
		val, err := strconv.Atoi(string(tt))
		if err != nil {
			return Null_
		}
		return Int(val)
	default:
		return Null_
	}
}

func float(args []OBJECT) OBJECT {
	switch tt := args[0].(type) {
	case Int:
		return Float(tt)
	case Float:
		return tt
	case *Bool:
		if tt == False_ {
			return Float(0)
		} else {
			return Float(1)
		}
	case *Null:
		return Float(0)
	case String:
		val, err := strconv.ParseFloat(string(tt), 64)
		if err != nil {
			return Null_
		}
		return Float(val)
	default:
		return Null_
	}
}
