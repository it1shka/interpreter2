package evaluator

import (
	"fmt"
	. "interpreter2/parser"
)

type Evaluator struct {
	scope *Scope
}

func NewEvaluator() *Evaluator {
	evaluator := new(Evaluator)
	evaluator.scope = CreateScope()
	evaluator._init()
	return evaluator
}

func (ev *Evaluator) _init() {
	ev.setBuiltins([]BuiltinRule{
		{BuiltinFunction(echo), "echo"},
	})
}

type BuiltinRule struct {
	fn         OBJECT
	identifier string
}

func (ev *Evaluator) setBuiltins(rules []BuiltinRule) {
	for _, rule := range rules {
		ev.scope.Init(rule.identifier, rule.fn)
	}
}

func (ev *Evaluator) pushScope() {
	ev.scope = ev.scope.Child()
}

func (ev *Evaluator) popScope() {
	ev.scope = ev.scope.prev
}

func (ev *Evaluator) Eval(ast STMT_LIST) {
	cb := ev.evalStmtList(ast)
	if cb != nil {
		panic("Raised unexpected callback")
	}
}

func (ev *Evaluator) evalStmtList(stmts STMT_LIST) Callback {
	for _, stmt := range stmts {
		switch t := stmt.(type) {

		case *IF_STMT:
			condition := ev.evalExpr(t.Condition)
			if checkCondition(condition) {
				ev.pushScope()
				cb := ev.evalStmtList(t.Then)
				ev.popScope()
				if cb != nil {
					return cb
				}
			} else if t.Els != nil {
				ev.pushScope()
				cb := ev.evalStmtList(t.Els)
				ev.popScope()
				if cb != nil {
					return cb
				}
			}

		case *FOR_STMT:
			for {
				condition := ev.evalExpr(t.Condition)
				if !checkCondition(condition) {
					break
				}
				ev.pushScope()
				cb := ev.evalStmtList(t.Body)
				ev.popScope()
				if _, ok := cb.(Break); ok {
					break
				}
				if cb, ok := cb.(Return); ok {
					return cb
				}
			}

		case *BREAK_STMT:
			return Break{}

		case *CONTINUE_STMT:
			return Continue{}

		case *RETURN_STMT:
			obj := ev.evalExpr(t.Expression)
			return Return{obj}

		case *LET_STMT:
			obj := ev.evalExpr(t.Expression)
			ev.scope.Init(t.Identifier, obj)

		case *EXPR_STMT:
			ev.evalExpr(t.Expression)
		}
	}
	return nil
}

func (ev *Evaluator) evalExpr(expression EXPRESSION) OBJECT {
	switch expr := expression.(type) {

	// Literals
	case *INT_LITERAL_EXPR:
		return Int(expr.Value)

	case *FLOAT_LITERAL_EXPR:
		return Float(expr.Value)

	case *BOOL_LITERAL_EXPR:
		return resolveBool(expr.Value)

	case *STR_LITERAL_EXPR:
		return String(expr.Value)

	case *ARRAY_LITERAL_EXPR:
		arr := make(Array, len(expr.Array))
		for ind, val := range expr.Array {
			arr[ind] = ev.evalExpr(val)
		}
		return arr

	case *FUNCTION_LITERAL_EXPR:
		return &Function{expr.Arguments, expr.Body, ev.scope}

	case *NULL_LITERAL_EXPR:
		return Null_

	// complex expressions

	case *VARIABLE_EXPR:
		return ev.scope.Get(expr.Identifier)

	case *BINARY_EXPR:
		right := ev.evalExpr(expr.Right)
		if expr.Operator == "=" {
			left, ok := expr.Left.(*VARIABLE_EXPR)
			if !ok {
				panic("Expected identifier in assignment")
			}
			ev.scope.SetOrInit(left.Identifier, right)
			return right
		} else {
			left := ev.evalExpr(expr.Left)
			return resolveBinaryOp(left, right, expr.Operator)
		}

	case *UNARY_EXPR:
		obj := ev.evalExpr(expr.Expression)
		return resolveUnaryOp(obj, expr.Operator)

	case *FUNCTION_CALL_EXPR:
		callArgs := make([]OBJECT, len(expr.Arguments))
		for ind, val := range expr.Arguments {
			callArgs[ind] = ev.evalExpr(val)
		}

		functionalObj := ev.evalExpr(expr.Function)

		switch fn := functionalObj.(type) {
		case BuiltinFunction:
			return fn(callArgs)
		case *Function:
			return ev.runFunction(fn, callArgs)
		default:
			panic(fmt.Sprintf("Object %s is not callable", fn.ToString()))
		}

	case *GET_OPERATOR_EXPR:
		fromObj := ev.evalExpr(expr.From)
		indexObj := ev.evalExpr(expr.Index)
		if fromObj.getType() == ARRAY_TYPE && indexObj.getType() == INT_TYPE {
			return fromObj.(Array)[indexObj.(Int)]
		} else {
			panic(fmt.Sprintf("Expected array and index, found %s and %s",
				fromObj.ToString(), indexObj.ToString()))
		}

	}

	panic("undefined behaviour")
}

func (ev *Evaluator) runFunction(fn *Function, args []OBJECT) OBJECT {
	backup := ev.scope
	defer func() {
		ev.scope = backup
	}()

	ev.scope = fn.OuterScope.Child()
	for i := 0; i < min(len(fn.Arguments), len(args)); i++ {
		ev.scope.Init(fn.Arguments[i], args[i])
	}

	cb := ev.evalStmtList(fn.Body)
	if cb != nil {
		if retCb, ok := cb.(Return); ok {
			return retCb.Value
		} else {
			panic("Unexpected callback in function")
		}
	}
	return Null_
}
