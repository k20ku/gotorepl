package gotorepl

import (
	"bytes"
	"fmt"
	"text/scanner"
)

const (
	EOF      rune = scanner.EOF
	LPAREN   rune = '('
	RPAREN   rune = ')'
	LBRACE   rune = '{'
	RBRACE   rune = '}'
	LBRACKET rune = '['
	RBRACKET rune = ']'
)

type Lexer struct {
	input    string
	scanner  *scanner.Scanner
	curToken rune
}

func NewLexer(input string) *Lexer {
	reader := bytes.NewBufferString(input)
	sc := new(scanner.Scanner)
	sc.Init(reader)
	sc.Error = func(_ *scanner.Scanner, msg string) {
		fmt.Printf("scanner: %s", msg)
	}
	return &Lexer{input: input, scanner: sc}
}

func (l *Lexer) NextToken() rune {
	return l.scanner.Scan()
}
