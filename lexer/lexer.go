package lexer

import (
	"github.com/zhulik/monkey/tokens"
)

type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
}

func New(input string) *Lexer {
	return &Lexer{input: input}
}

func (l *Lexer) NextToken() tokens.Token {
	l.readChar()

	var tok tokens.Token

	switch l.ch {
	case '=':
		tok = tokens.Token{Type: tokens.ASSIGN, Literal: string(l.ch)}
	case '+':
		tok = tokens.Token{Type: tokens.PLUS, Literal: string(l.ch)}
	case '(':
		tok = tokens.Token{Type: tokens.LPAREN, Literal: string(l.ch)}
	case ')':
		tok = tokens.Token{Type: tokens.RPAREN, Literal: string(l.ch)}
	case '{':
		tok = tokens.Token{Type: tokens.LBRACE, Literal: string(l.ch)}
	case '}':
		tok = tokens.Token{Type: tokens.RBRACE, Literal: string(l.ch)}
	case ',':
		tok = tokens.Token{Type: tokens.COMMA, Literal: string(l.ch)}
	case ';':
		tok = tokens.Token{Type: tokens.SEMICOLON, Literal: string(l.ch)}
	case 0:
		tok = tokens.Token{Type: tokens.EOF}
	}

	return tok
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}

	l.position = l.readPosition
	l.readPosition++
}
