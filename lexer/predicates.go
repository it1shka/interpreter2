package lexer

import (
	"strings"
	"unicode"
)

const allowedSpecialSymbols = "$@_?"

func isWhitespace(char rune) bool {
	return unicode.IsSpace(char)
}

func ifNot(char rune) func(rune) bool {
	return func(charPeeked rune) bool {
		return charPeeked != char
	}
}

/*
func notNewLine(char rune) bool {
	return char != '\n'
}
*/

var notNewLine = ifNot('\n')

func isNumStart(char rune) bool {
	return unicode.IsDigit(char)
}

func isWordStart(char rune) bool {
	return unicode.IsLetter(char) || strings.ContainsRune(allowedSpecialSymbols, char)
}

func isWordLetter(char rune) bool {
	return isNumStart(char) || isWordStart(char)
}
