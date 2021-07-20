package lexer

import "fmt"

type TokenType int

const (
	EOF_TOK TokenType = iota

	ID_TOK
	KW_TOK
	PUNC_TOK
	OP_TOK

	INT_TOK
	FLOAT_TOK
	STR_TOK
)

func (t TokenType) Format() string {
	switch t {
	case EOF_TOK:
		return "end of file"
	case ID_TOK:
		return "identifier token"
	case KW_TOK:
		return "keyword token"
	case PUNC_TOK:
		return "punctuation token"
	case OP_TOK:
		return "operator token"
	case INT_TOK:
		return "integer literal token"
	case FLOAT_TOK:
		return "floating literal token"
	case STR_TOK:
		return "string literal token"
	}
	return ""
}

type Token struct {
	Type    TokenType
	Literal string
}

func (token *Token) Format() string {
	return fmt.Sprintf("\"%s\" of type %s",
		token.Literal, token.Type.Format())
}
