package lexer

func RunLexer(code string) []*Token {
	tokens := make([]*Token, 0)
	lexer := NewLexer(code)
	for !lexer.Eof() {
		tokens = append(tokens, lexer.NextToken())
	}
	return tokens
}
