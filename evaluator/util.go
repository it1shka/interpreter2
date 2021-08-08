package evaluator

// import "encoding/json"

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
	case Array:
		return len(obj) > 0
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
func addf(a, b float64) float64 { return a + b }

func resolveBinaryOp(left, right OBJECT, operator string) OBJECT {

	switch operator {
	case "|":
		return resolveBool(checkCondition(left) || checkCondition(right))
	case "&":
		return resolveBool(checkCondition(left) && checkCondition(right))
	case "==":
		return resolveBool(left == right)
	case "!=":
		return resolveBool(left != right)
	case "+":
		if obj, ok := TryNumberOperation(left, right, addf); ok {
			return obj
		}

		if left.getType() == STRING_TYPE && right.getType() == STRING_TYPE {
			return left.(String) + right.(String)
		}

		if left.getType() == ARRAY_TYPE {
			if right.getType() == ARRAY_TYPE {
				return append(left.(Array), right.(Array)...)
			} else {
				return append(left.(Array), right)
			}
		}

		if right.getType() == ARRAY_TYPE {
			if left.getType() == ARRAY_TYPE {
				return append(left.(Array), right.(Array)...)
			} else {
				return append(Array{left}, right.(Array)...)
			}
		}

	}
	panic("Unfinished operator")
}

func resolveUnaryOp(obj OBJECT, operator string) OBJECT {
	panic("Unfinished operator")
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
