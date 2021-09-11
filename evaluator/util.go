package evaluator

// import "encoding/json"
import "strings"

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func resolveBool(v bool) *Bool {
	if v {
		return True_
	}
	return False_
}

func checkCondition(object OBJECT) bool {
	switch obj := object.(type) {
	case *Bool:
		return bool(*obj)
	case Int:
		return obj != 0
	case Float:
		return obj != 0
	case *Null:
		return false
	case String:
		return len(obj) > 0
	case *Array:
		return len(*obj) > 0
	default:
		return true
	}
}

func areNumbers(a, b OBJECT) bool {
	return (a.getType() == INT_TYPE || a.getType() == FLOAT_TYPE) &&
		(b.getType() == INT_TYPE || b.getType() == FLOAT_TYPE)
}

func ToFloat64(a OBJECT) (out float64) {
	switch a_ := a.(type) {
	case Float:
		out = float64(a_)
	case Int:
		out = float64(a_)
	}
	return
}

type NumberOpFunc func(a, b float64) float64

type NumberCompFunc func(a, b float64) bool

func NumberOperation(a, b OBJECT, fn NumberOpFunc) OBJECT {
	af, bf := ToFloat64(a), ToFloat64(b)
	result := fn(af, bf)
	if a.getType() == INT_TYPE && b.getType() == INT_TYPE {
		return Int(result)
	}
	return Float(result)
}

func TryNumberOperation(a, b OBJECT, fn NumberOpFunc) (OBJECT, bool) {
	if areNumbers(a, b) {
		return NumberOperation(a, b, fn), true
	}
	return Null_, false
}

func TryNumberCompare(a, b OBJECT, fn NumberCompFunc) (OBJECT, bool) {
	if areNumbers(a, b) {
		result := fn(ToFloat64(a), ToFloat64(b))
		return resolveBool(result), true
	}
	return Null_, false
}

// number functions...

func eqf(a, b float64) bool { return a == b }
func lsf(a, b float64) bool { return a < b }
func gtf(a, b float64) bool { return a > b }

// func nef(a, b float64) bool { return a != b }
func addf(a, b float64) float64 { return a + b }
func subf(a, b float64) float64 { return a - b }
func mulf(a, b float64) float64 { return a * b }
func divf(a, b float64) float64 { return a / b }
func modf(a, b float64) float64 { return float64(int(a) % int(b)) }

func resolveBinaryOp(left, right OBJECT, operator string) OBJECT {

	switch operator {
	case "|":
		return resolveBool(checkCondition(left) || checkCondition(right))
	case "&":
		return resolveBool(checkCondition(left) && checkCondition(right))
	case "==":
		if obj, ok := TryNumberCompare(left, right, eqf); ok {
			return obj
		}
		if left.getType() != right.getType() {
			return False_
		}
		if left.getType() == ARRAY_TYPE {
			left_ := left.(*Array)
			right_ := right.(*Array)
			if len(*left_) != len(*right_) {
				return False_
			}
			for i := 0; i < len(*left_); i++ {
				res := resolveBinaryOp((*left_)[i], (*right_)[i], "==")
				if res == False_ {
					return False_
				}
			}
			return True_
		}
		return resolveBool(left == right)
	case "!=":
		res := resolveBinaryOp(left, right, "==")
		return resolveUnaryOp(res, "!")
	case "<":
		if obj, ok := TryNumberCompare(left, right, lsf); ok {
			return obj
		}
		if left.getType() == STRING_TYPE && left.getType() == STRING_TYPE {
			return resolveBool(left.(String) < right.(String))
		}
		if left.getType() == ARRAY_TYPE && right.getType() == ARRAY_TYPE {
			left_ := left.(*Array)
			right_ := right.(*Array)
			if len(*left_) != len(*right_) {
				return resolveBool(len(*left_) < len(*right_))
			}
			return False_
		}
		return False_
	case ">":
		if obj, ok := TryNumberCompare(left, right, gtf); ok {
			return obj
		}
		if left.getType() == STRING_TYPE && left.getType() == STRING_TYPE {
			return resolveBool(left.(String) > right.(String))
		}
		if left.getType() == ARRAY_TYPE && right.getType() == ARRAY_TYPE {
			left_ := left.(*Array)
			right_ := right.(*Array)
			if len(*left_) != len(*right_) {
				return resolveBool(len(*left_) > len(*right_))
			}
			return False_
		}
		return False_
	case "<=":
		res := resolveBinaryOp(left, right, ">")
		return resolveUnaryOp(res, "!")
	case ">=":
		res := resolveBinaryOp(left, right, "<")
		return resolveUnaryOp(res, "!")

	case "+":
		if obj, ok := TryNumberOperation(left, right, addf); ok {
			return obj
		}

		if left.getType() == STRING_TYPE || right.getType() == STRING_TYPE {
			return String(left.ToString() + right.ToString())
		}

		if left.getType() == ARRAY_TYPE {
			left_ := left.(*Array)
			if right.getType() == ARRAY_TYPE {
				right_ := right.(*Array)
				narr := append(*left_, *right_...)
				return &narr
			} else {
				narr := append(*left_, right)
				return &narr
			}
		}

		if right.getType() == ARRAY_TYPE {
			right_ := right.(*Array)
			if left.getType() == ARRAY_TYPE {
				left_ := left.(*Array)
				narr := append(*right_, *left_...)
				return &narr
			} else {
				narr := append(*right_, left)
				return &narr
			}
		}

		return Null_
	case "-":
		if obj, ok := TryNumberOperation(left, right, subf); ok {
			return obj
		}
		if left.getType() == ARRAY_TYPE {
			left_ := left.(*Array)
			for ind, each := range *left_ {
				if resolveBinaryOp(each, right, "==") == True_ {
					narr := append((*left_)[:ind], (*left_)[ind+1:]...)
					return &narr
				}
			}
			return left
		}
		if left.getType() == STRING_TYPE && right.getType() == STRING_TYPE {
			res := strings.Replace(string(left.(String)), string(right.(String)), "", 1)
			return String(res)
		}

		return Null_
	case "*":
		if obj, ok := TryNumberOperation(left, right, mulf); ok {
			return obj
		}
		if left.getType() == STRING_TYPE && right.getType() == INT_TYPE {
			res := strings.Repeat(string(left.(String)), int(right.(Int)))
			return String(res)
		}
		if right.getType() == STRING_TYPE && left.getType() == INT_TYPE {
			res := strings.Repeat(string(right.(String)), int(left.(Int)))
			return String(res)
		}

		if left.getType() == ARRAY_TYPE && right.getType() == INT_TYPE {
			left_ := left.(*Array)
			if right.(Int) < 0 {
				return Null_
			}
			res := make(Array, 0)
			for i := 0; i < int(right.(Int)); i++ {
				res = append(res, *left_...)
			}
			return &res
		}
		if right.getType() == ARRAY_TYPE && left.getType() == INT_TYPE {
			right_ := right.(*Array)
			if left.(Int) < 0 {
				return Null_
			}
			res := make(Array, 0)
			for i := 0; i < int(left.(Int)); i++ {
				res = append(res, *right_...)
			}
			return &res
		}

		return Null_
	case "/":
		if obj, ok := TryNumberOperation(left, right, divf); ok {
			return obj
		}
		return Null_
	case "%":
		if obj, ok := TryNumberOperation(left, right, modf); ok {
			return obj
		}
		return Null_
	}

	panic("Undefined behaviour")
}

func resolveUnaryOp(obj OBJECT, operator string) OBJECT {
	switch operator {
	case "!":
		return resolveBool(!checkCondition(obj))
	case "-":
		return resolveBinaryOp(obj, Int(-1), "*")
	}

	panic("Undefined behaviour")
}

/*
func FormatStruct(structure interface{}) string {
	s, err := json.MarshalIndent(structure, "", "\t")
	if err != nil {
		panic("Failed to format structure")
	}
	return string(s)
}
*/
