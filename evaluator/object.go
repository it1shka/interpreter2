package evaluator

import (
	"fmt"
	"interpreter2/parser"
	"strings"
)

type OBJECT interface {
	_obj()
	getType() string
	ToString() string
}

type Int int

type Float float64

type Bool bool

type Null struct{}

type String string

type Array []OBJECT

type Function struct {
	Arguments  []string
	Body       parser.STMT_LIST
	OuterScope *Scope
}

// for builtins
type BuiltinFunction func([]OBJECT) OBJECT

// binding to OBJECT interface

func (s Int) _obj()             {}
func (s Float) _obj()           {}
func (s *Bool) _obj()           {}
func (s *Null) _obj()           {}
func (s String) _obj()          {}
func (s Array) _obj()           {}
func (s *Function) _obj()       {}
func (s BuiltinFunction) _obj() {}

func (s Int) getType() string             { return "Int" }
func (s Float) getType() string           { return "Float" }
func (s *Bool) getType() string           { return "Bool" }
func (s *Null) getType() string           { return "Null" }
func (s String) getType() string          { return "String" }
func (s Array) getType() string           { return "Array" }
func (s *Function) getType() string       { return "Function" }
func (s BuiltinFunction) getType() string { return "BuiltIn" }

func (s Int) ToString() string    { return fmt.Sprintf("%d", s) }
func (s Float) ToString() string  { return fmt.Sprintf("%f", s) }
func (s *Bool) ToString() string  { return fmt.Sprintf("%t", *s) }
func (s *Null) ToString() string  { return "null" }
func (s String) ToString() string { return string(s) }
func (s Array) ToString() string {
	str := make([]string, 0)
	for _, object := range s {
		str = append(str, object.ToString())
	}
	return "[" + strings.Join(str, ", ") + "]"
}
func (s *Function) ToString() string {
	return "functional object"
}
func (s BuiltinFunction) ToString() string {
	return "builtin object"
}

// global variables

var true_ = Bool(true)
var false_ = Bool(false)

var True_ = &true_
var False_ = &false_
var Null_ = &Null{}
