package lexer

import (
	"fmt"
	"strconv"
	"strings"
)

type Lexer struct {
	source   []rune
	position int
	current  *Token
	char     rune
}

func NewLexer(source string) *Lexer {
	lexer := new(Lexer)
	lexer.source = []rune(source)
	lexer.position = 0
	return lexer
}

func (lexer *Lexer) peek() rune {
	if lexer.position >= len(lexer.source) {
		return rune(0)
	}
	return lexer.source[lexer.position]
}

func (lexer *Lexer) next() rune {
	character := lexer.peek()
	lexer.position++
	return character
}

func (lexer *Lexer) isEof() bool {
	return lexer.peek() == rune(0)
}

func (lexer *Lexer) readWhile(predicate func(rune) bool) string {
	var stringBuilder strings.Builder
	for !lexer.isEof() && predicate(lexer.peek()) {
		stringBuilder.WriteRune(lexer.next())
	}
	return stringBuilder.String()
}

func (lexer *Lexer) eatWhitespace() {
	lexer.readWhile(isWhitespace)
}

func (lexer *Lexer) eatComment() {
	lexer.readWhile(notNewLine)
}

func (lexer *Lexer) readToken() Token {
	lexer.eatWhitespace()
	if lexer.isEof() {
		return Token{EOF_TOK, ""}
	}
	lexer.char = lexer.next()
	switch lexer.char {
	case '#':
		lexer.eatComment()
		return lexer.readToken()
	case '(', ')', '{', '}', ',', ';', '[', ']':
		return Token{PUNC_TOK, string(lexer.char)}
	case '+', '-', '*', '/', '%', '&', '|':
		return Token{OP_TOK, string(lexer.char)}
	case '=', '!', '>', '<':
		if lexer.peek() == '=' {
			return Token{OP_TOK, fmt.Sprintf("%c%c", lexer.char, lexer.next())}
		}
		return Token{OP_TOK, string(lexer.char)}
	case '\'', '"':
		return lexer.readStringLiteral()
	}
	if isNumStart(lexer.char) {
		return lexer.readNumberLiteral()
	}
	if isWordStart(lexer.char) {
		return lexer.readWord()
	}
	panic(fmt.Sprintf("Unexpected character: %c", lexer.char))
}

func (lexer *Lexer) readStringLiteral() Token {
	quoteType := lexer.char
	stringValue := lexer.readWhile(ifNot(quoteType))
	enclosingQuote := lexer.next()
	if enclosingQuote != quoteType {
		panic(fmt.Sprintf("Expected string closing quote"))
	}

	switch quoteType {
	case '\'':
		// raw string
		stringValue = fmt.Sprintf("`%s`", stringValue)
	case '"':
		// string with special characters like \n \t ...
		stringValue = fmt.Sprintf("\"%s\"", stringValue)
	}

	unquotedValue, err := strconv.Unquote(stringValue)
	if err != nil {
		panic("Invalid string syntax")
	}
	return Token{STR_TOK, unquotedValue}
}

func (lexer *Lexer) readNumberLiteral() Token {
	readNumbers := func() string {
		return lexer.readWhile(isNumStart)
	}

	numberValue := string(lexer.char) + readNumbers()
	if lexer.peek() == '.' {
		numberValue += string(lexer.next()) + readNumbers()
		return Token{FLOAT_TOK, numberValue}
	}
	return Token{INT_TOK, numberValue}
}

func (lexer *Lexer) readWord() Token {
	wordValue := string(lexer.char) + lexer.readWhile(isWordLetter)
	switch wordValue {
	case "if", "for", "break",
		"else", "continue", "return",
		"let", "fn", "true", "false", "null":
		return Token{KW_TOK, wordValue}
	default:
		return Token{ID_TOK, wordValue}
	}
}

func (lexer *Lexer) PeekToken() *Token {
	if lexer.current == nil {
		token := lexer.readToken()
		lexer.current = &token
	}
	return lexer.current
}

func (lexer *Lexer) NextToken() *Token {
	token := lexer.PeekToken()
	lexer.current = nil
	return token
}

func (lexer *Lexer) Eof() bool {
	return lexer.PeekToken().Type == EOF_TOK
}

func (lexer *Lexer) NextTokenIf(expectedLiteral string) bool {
	if lexer.PeekToken().Literal == expectedLiteral {
		lexer.NextToken()
		return true
	}
	return false
}

func (lexer *Lexer) ExpectToken(expectedLiteral string) {
	if lexer.NextToken().Literal != expectedLiteral {
		errStr := fmt.Sprintf("Expected \"%s\" literal", expectedLiteral)
		panic(errStr)
	}
}
