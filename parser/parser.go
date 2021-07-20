package parser

import (
	"fmt"
	. "interpreter2/lexer"
	"strconv"
)

type Parser struct {
	buff            *Lexer
	parseExpression func() EXPRESSION
}

func NewParser(code string) *Parser {
	parser := &Parser{NewLexer(code), nil}
	parser._init()
	return parser
}

func (parser *Parser) _init() {
	parser.parseExpression = parser.combineRules(OperatorRules{
		{parser.parseRightOps, []string{"="}},
		{parser.parseLeftOps, []string{"|"}},
		{parser.parseLeftOps, []string{"&"}},
		{parser.parseLeftOps, []string{"==", "!="}},
		{parser.parseLeftOps, []string{"<", ">", "<=", ">="}},
		{parser.parseLeftOps, []string{"-", "+", "*"}},
		{parser.parseLeftOps, []string{"/", "%"}},
		{parser.parseUnaryOps, []string{"-", "!"}},
		{parser.parsePostfixOps, []string{"(", "["}},
	})
}

func (parser *Parser) ParseProgram() STMT_LIST {
	if parser.buff.Eof() {
		return []STATEMENT{}
	}
	statement := parser.parseStatement()
	return append([]STATEMENT{statement}, parser.ParseProgram()...)
}

func (parser *Parser) parseStatement() STATEMENT {
	switch parser.buff.PeekToken().Literal {
	case "if":
		parser.buff.NextToken()
		condition := parser.parseExpression()
		body := parser.parseStatementList()
		var elseBody STMT_LIST
		if parser.buff.NextTokenIf("else") {
			elseBody = parser.parseStatementList()
		}
		return &IF_STMT{condition, body, elseBody}
	case "for":
		parser.buff.NextToken()
		condition := parser.parseExpression()
		body := parser.parseStatementList()
		return &FOR_STMT{condition, body}
	case "break":
		parser.buff.NextToken()
		parser.buff.ExpectToken(";")
		return &BREAK_STMT{}
	case "continue":
		parser.buff.NextToken()
		parser.buff.ExpectToken(";")
		return &CONTINUE_STMT{}
	case "return":
		parser.buff.NextToken()
		expression := parser.parseExpression()
		parser.buff.ExpectToken(";")
		return &RETURN_STMT{expression}
	case "let":
		parser.buff.NextToken()
		if parser.buff.PeekToken().Type != ID_TOK {
			panic("Expected identifier")
		}
		identifier := parser.buff.NextToken().Literal
		var expression EXPRESSION
		if parser.buff.NextTokenIf("=") {
			expression = parser.parseExpression()
		}
		parser.buff.ExpectToken(";")
		return &LET_STMT{identifier, expression}
	default:
		expression := parser.parseExpression()
		parser.buff.ExpectToken(";")
		return &EXPR_STMT{expression}
	}
}

type OperatorRules []struct {
	fn        func(func() EXPRESSION, []string) EXPRESSION
	operators []string
}

func (parser *Parser) combineRules(rules OperatorRules) func() EXPRESSION {
	if len(rules) == 0 {
		return parser.parseValue
	}
	currentRule := rules[0]
	return func() EXPRESSION {
		return currentRule.fn(parser.combineRules(rules[1:]), currentRule.operators)
	}
}

func (parser *Parser) parseExpressionWithDelimiter(del string, end string) []EXPRESSION {
	expressions := make([]EXPRESSION, 0)
	for {
		if parser.buff.PeekToken().Literal == end {
			break
		}
		expression := parser.parseExpression()
		expressions = append(expressions, expression)
		if !parser.buff.NextTokenIf(del) {
			break
		}
	}
	parser.buff.ExpectToken(end)
	return expressions
}

func (parser *Parser) parseIdentifiersWithDelimiter(del string, end string) []string {
	identifiers := make([]string, 0)
	for {
		if parser.buff.PeekToken().Literal == end {
			break
		}
		identifier := parser.buff.NextToken()
		if identifier.Type != ID_TOK {
			panic("Identifier expected")
		}
		identifiers = append(identifiers, identifier.Literal)
		if !parser.buff.NextTokenIf(del) {
			break
		}
	}
	return identifiers
}

func (parser *Parser) parseValue() EXPRESSION {
	current := parser.buff.NextToken()
	t_type, t_literal := current.Type, current.Literal
	switch t_literal {
	case "(":
		expression := parser.parseExpression()
		parser.buff.ExpectToken(")")
		return expression
	case "true":
		return &BOOL_LITERAL_EXPR{true}
	case "false":
		return &BOOL_LITERAL_EXPR{false}
	case "null":
		return &NULL_LITERAL_EXPR{}
	case "[":
		expressions := parser.parseExpressionWithDelimiter(",", "]")
		return &ARRAY_LITERAL_EXPR{expressions}
	case "fn":
		arguments := parser.parseIdentifiersWithDelimiter(",", "{")
		body := parser.parseStatementList()
		return &FUNCTION_LITERAL_EXPR{arguments, body}
	}

	switch t_type {
	case ID_TOK:
		return &VARIABLE_EXPR{t_literal}
	case INT_TOK:
		val, err := strconv.Atoi(t_literal)
		if err != nil {
			panic("Failed to parse integer")
		}
		return &INT_LITERAL_EXPR{val}
	case FLOAT_TOK:
		val, err := strconv.ParseFloat(t_literal, 64)
		if err != nil {
			panic("Failed to parse floating")
		}
		return &FLOAT_LITERAL_EXPR{val}
	case STR_TOK:
		return &STR_LITERAL_EXPR{t_literal}
	default:
		errStr := fmt.Sprintf("Expected int, float, string, array or function, found: %s",
			current.Format())
		panic(errStr)
	}
}

func (parser *Parser) parseUnaryOps(fn func() EXPRESSION, operators []string) EXPRESSION {
	if Includes(parser.buff.PeekToken().Literal, operators) {
		operator := parser.buff.NextToken().Literal
		right := parser.parseUnaryOps(fn, operators)
		return &UNARY_EXPR{operator, right}
	}
	return fn()
}

func (parser *Parser) parsePostfixOps(fn func() EXPRESSION, operators []string) EXPRESSION {
	left := fn()
	for Includes(parser.buff.PeekToken().Literal, operators) {
		left = parser.parsePostifxOp(left)
	}
	return left
}

func (parser *Parser) parsePostifxOp(expr EXPRESSION) EXPRESSION {
	switch parser.buff.NextToken().Literal {
	case "[":
		index := parser.parseExpression()
		parser.buff.ExpectToken("]")
		return &GET_OPERATOR_EXPR{expr, index}
	case "(":
		expressions := parser.parseExpressionWithDelimiter(",", ")")
		return &FUNCTION_CALL_EXPR{expr, expressions}
	default:
		panic("Expected postfix operator")
	}
}

func (parser *Parser) parseLeftOps(fn func() EXPRESSION, operators []string) EXPRESSION {
	result := fn()
	for Includes(parser.buff.PeekToken().Literal, operators) {
		operator := parser.buff.NextToken().Literal
		right := fn()
		result = &BINARY_EXPR{operator, result, right}
	}
	return result
}

func (parser *Parser) parseRightOps(fn func() EXPRESSION, operators []string) EXPRESSION {
	result := fn()
	if Includes(parser.buff.PeekToken().Literal, operators) {
		operator := parser.buff.NextToken().Literal
		right := parser.parseRightOps(fn, operators)
		result = &BINARY_EXPR{operator, result, right}
	}
	return result
}

func (parser *Parser) parseStatementList() STMT_LIST {
	parser.buff.ExpectToken("{")
	statements := make([]STATEMENT, 0)
	for {
		if parser.buff.NextTokenIf("}") {
			break
		}
		statement := parser.parseStatement()
		statements = append(statements, statement)
	}
	return statements
}
