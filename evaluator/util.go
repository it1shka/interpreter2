package evaluator

// import "encoding/json"

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
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

func resolveBinaryOp(left, right OBJECT, operator string) OBJECT {

	switch operator {
	case "+":
		if left.getType() == INT_TYPE && right.getType() == INT_TYPE {
			return Int(left.(Int) + right.(Int))
		}

		if left.getType() == FLOAT_TYPE && right.getType() == FLOAT_TYPE {
			return Float(left.(Float) + right.(Float))
		}

		if left.getType() == FLOAT_TYPE && right.getType() == INT_TYPE {
			return Float(left.(Float) + Float(right.(Int)))
		}

		if left.getType() == INT_TYPE && right.getType() == FLOAT_TYPE {
			return Float(Float(left.(Int)) + right.(Float))
		}

		panic("Unfinished operator")
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
