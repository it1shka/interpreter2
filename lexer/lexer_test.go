package lexer

import (
	"fmt"
	"strings"
	"testing"
)

func CheckTokensEquality(slice1, slice2 []*Token) bool {
	if len(slice1) != len(slice2) {
		fmt.Println("slices are of different sizes")
		fmt.Println("One is: ", FormatTokens(slice1))
		fmt.Println("Another is: ", FormatTokens(slice2))
		return false
	}
	for i := range slice1 {
		tok1, tok2 := slice1[i], slice2[i]
		if tok1.Literal != tok2.Literal || tok1.Type != tok2.Type {
			fmt.Printf("Failed to compare: %s AND %s", tok1.Format(), tok2.Format())
			return false
		}
	}
	return true
}

func FormatTokens(tokens []*Token) string {
	var sb strings.Builder
	for _, val := range tokens {
		sb.WriteString(val.Format())
		sb.WriteRune('\n')
	}
	return sb.String()
}

func TestLexer(t *testing.T) {
	fmt.Println("TESTING NORMAL CASES: ")

	tables := []struct {
		code   string
		result []*Token
	}{
		{
			`let greetings = "Hello world!"`,
			[]*Token{
				{KW_TOK, "let"},
				{ID_TOK, "greetings"},
				{OP_TOK, "="},
				{STR_TOK, "Hello world!"},
			},
		},
		{
			"\n\n\n        \t\t\r\r\r    \n#bruh i love coding\necho(\"Hello world!\")",
			[]*Token{
				{ID_TOK, "echo"},
				{PUNC_TOK, "("},
				{STR_TOK, "Hello world!"},
				{PUNC_TOK, ")"},
			},
		},
		{
			`
			#this is example program :)
			let counter = 1;
			for counter < 100 {
				echo(counter);
				counter = counter + 1;
			}
			`,
			[]*Token{
				{KW_TOK, "let"},
				{ID_TOK, "counter"},
				{OP_TOK, "="},
				{INT_TOK, "1"},
				{PUNC_TOK, ";"},
				{KW_TOK, "for"},
				{ID_TOK, "counter"},
				{OP_TOK, "<"},
				{INT_TOK, "100"},
				{PUNC_TOK, "{"},
				{ID_TOK, "echo"},
				{PUNC_TOK, "("},
				{ID_TOK, "counter"},
				{PUNC_TOK, ")"},
				{PUNC_TOK, ";"},
				{ID_TOK, "counter"},
				{OP_TOK, "="},
				{ID_TOK, "counter"},
				{OP_TOK, "+"},
				{INT_TOK, "1"},
				{PUNC_TOK, ";"},
				{PUNC_TOK, "}"},
			},
		},
		{
			"10.5, 10.6, 10.666666666666666, 101, 666, bruh bruh101 bruh___?, $@$@$@",
			[]*Token{
				{FLOAT_TOK, "10.5"},
				{PUNC_TOK, ","},
				{FLOAT_TOK, "10.6"},
				{PUNC_TOK, ","},
				{FLOAT_TOK, "10.666666666666666"},
				{PUNC_TOK, ","},
				{INT_TOK, "101"},
				{PUNC_TOK, ","},
				{INT_TOK, "666"},
				{PUNC_TOK, ","},
				{ID_TOK, "bruh"},
				{ID_TOK, "bruh101"},
				{ID_TOK, "bruh___?"},
				{PUNC_TOK, ","},
				{ID_TOK, "$@$@$@"},
			},
		},
		{
			"if \nfor \t\nbreak\n\r\n else continue        return let fn",
			[]*Token{
				{KW_TOK, "if"},
				{KW_TOK, "for"},
				{KW_TOK, "break"},
				{KW_TOK, "else"},
				{KW_TOK, "continue"},
				{KW_TOK, "return"},
				{KW_TOK, "let"},
				{KW_TOK, "fn"},
			},
		},
		{
			"\"Hello,\\nworld!\" 'Hello,\nworld!'",
			[]*Token{
				{STR_TOK, "Hello,\nworld!"},
				{STR_TOK, "Hello,\nworld!"},
			},
		},
		{
			";{}()[],",
			[]*Token{
				{PUNC_TOK, ";"},
				{PUNC_TOK, "{"},
				{PUNC_TOK, "}"},
				{PUNC_TOK, "("},
				{PUNC_TOK, ")"},
				{PUNC_TOK, "["},
				{PUNC_TOK, "]"},
				{PUNC_TOK, ","},
			},
		},
		{
			"!-= == >>=> <=!= %",
			[]*Token{
				{OP_TOK, "!"},
				{OP_TOK, "-"},
				{OP_TOK, "="},
				{OP_TOK, "=="},
				{OP_TOK, ">"},
				{OP_TOK, ">="},
				{OP_TOK, ">"},
				{OP_TOK, "<="},
				{OP_TOK, "!="},
				{OP_TOK, "%"},
			},
		},
		{
			"true false null",
			[]*Token{
				{KW_TOK, "true"},
				{KW_TOK, "false"},
				{KW_TOK, "null"},
			},
		},
		{
			"let a = true; let b = false; echo(a & b);",
			[]*Token{
				{KW_TOK, "let"},
				{ID_TOK, "a"},
				{OP_TOK, "="},
				{KW_TOK, "true"},
				{PUNC_TOK, ";"},
				{KW_TOK, "let"},
				{ID_TOK, "b"},
				{OP_TOK, "="},
				{KW_TOK, "false"},
				{PUNC_TOK, ";"},
				{ID_TOK, "echo"},
				{PUNC_TOK, "("},
				{ID_TOK, "a"},
				{OP_TOK, "&"},
				{ID_TOK, "b"},
				{PUNC_TOK, ")"},
				{PUNC_TOK, ";"},
			},
		},
		{
			"echo(\"All tests are passed 😊😊😊\");",
			[]*Token{
				{ID_TOK, "echo"},
				{PUNC_TOK, "("},
				{STR_TOK, "All tests are passed 😊😊😊"},
				{PUNC_TOK, ")"},
				{PUNC_TOK, ";"},
			},
		},
	}

	for i, table := range tables {
		fmt.Println(i)
		fmt.Println(table.code)
		fmt.Println()
		tokens := RunLexer(table.code)
		fmt.Println(FormatTokens(tokens))
		fmt.Println()
		if !CheckTokensEquality(tokens, table.result) {
			t.Error("failed")
		}
	}
}

func TestPanics(t *testing.T) {
	tests := []string{
		"let a = .123",
		"let stringwithoutclosing = \"",
		"let some weird stuff = `",
		"another weird stuff = ^^^^^^^^^^",
		"wow this is smiley 😊😊😊",
		"this is a wave wow!!!!!!!!! ~~~~~~~~~~~~~~~",
	}

	runTest := func(test string) {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Panic and this is normal")
			}
		}()
		RunLexer(test)
		t.Errorf("Didnt panic but it had to. Fail.")
	}

	for _, test := range tests {
		runTest(test)
	}
}
