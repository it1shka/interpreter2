package parser

// statements

type STATEMENT interface {
	_stmt()
}

type STMT_LIST []STATEMENT

type IF_STMT struct {
	Condition EXPRESSION
	Then      STMT_LIST
	Els       STMT_LIST
}

type FOR_STMT struct {
	Condition EXPRESSION
	Body      STMT_LIST
}

type BREAK_STMT struct{}

type CONTINUE_STMT struct{}

type RETURN_STMT struct {
	Expression EXPRESSION
}

type LET_STMT struct {
	Identifier string
	Expression EXPRESSION
}

type EXPR_STMT struct {
	Expression EXPRESSION
}

// binding to STATEMENT interface
func (s *IF_STMT) _stmt()       {}
func (s *FOR_STMT) _stmt()      {}
func (s *BREAK_STMT) _stmt()    {}
func (s *CONTINUE_STMT) _stmt() {}
func (s *RETURN_STMT) _stmt()   {}
func (s *LET_STMT) _stmt()      {}
func (s *EXPR_STMT) _stmt()     {}

// expressions

type EXPRESSION interface {
	_expr()
}

type BINARY_EXPR struct {
	Operator    string
	Left, Right EXPRESSION
}

type UNARY_EXPR struct {
	Operator   string
	Expression EXPRESSION
}

type FUNCTION_CALL_EXPR struct {
	Function  EXPRESSION
	Arguments []EXPRESSION
}

type GET_OPERATOR_EXPR struct {
	From  EXPRESSION
	Index EXPRESSION
}

type VARIABLE_EXPR struct {
	Identifier string
}

// raw type expressions

type INT_LITERAL_EXPR struct {
	Value int
}

type FLOAT_LITERAL_EXPR struct {
	Value float64
}

type BOOL_LITERAL_EXPR struct {
	Value bool
}

type STR_LITERAL_EXPR struct {
	Value string
}

type ARRAY_LITERAL_EXPR struct {
	Array []EXPRESSION
}

type FUNCTION_LITERAL_EXPR struct {
	Arguments []string
	Body      STMT_LIST
}

type NULL_LITERAL_EXPR struct{}

// binding to expression interface
func (s *BINARY_EXPR) _expr()           {}
func (s *UNARY_EXPR) _expr()            {}
func (s *FUNCTION_CALL_EXPR) _expr()    {}
func (s *GET_OPERATOR_EXPR) _expr()     {}
func (s *VARIABLE_EXPR) _expr()         {}
func (s *INT_LITERAL_EXPR) _expr()      {}
func (s *FLOAT_LITERAL_EXPR) _expr()    {}
func (s *BOOL_LITERAL_EXPR) _expr()     {}
func (s *STR_LITERAL_EXPR) _expr()      {}
func (s *ARRAY_LITERAL_EXPR) _expr()    {}
func (s *FUNCTION_LITERAL_EXPR) _expr() {}
func (s *NULL_LITERAL_EXPR) _expr()     {}
