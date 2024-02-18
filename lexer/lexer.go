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
	l := &Lexer{input: input}
	l.readChar()

	return l
}

func (l *Lexer) NextToken() tokens.Token { //nolint:cyclop
	l.skipWhitespaces()

	var tok tokens.Token

	switch l.ch {
	case '=':
		tok = tokens.New(tokens.ASSIGN)
	case '+':
		tok = tokens.New(tokens.PLUS)
	case '(':
		tok = tokens.New(tokens.LPAREN)
	case ')':
		tok = tokens.New(tokens.RPAREN)
	case '{':
		tok = tokens.New(tokens.LBRACE)
	case '}':
		tok = tokens.New(tokens.RBRACE)
	case ',':
		tok = tokens.New(tokens.COMMA)
	case ';':
		tok = tokens.New(tokens.SEMICOLON)
	case 0:
		tok = tokens.New(tokens.EOF)
	default:
		switch {
		case isLetter(l.ch):
			literal := l.readIdentifier()

			tok.Type = tokens.TypeFromLiteral(literal)
			if tok.Type == tokens.IDENTIFIER {
				tok.Literal = literal
			}

			return tok
		case isDigit(l.ch):
			tok.Literal = l.readNumber()
			tok.Type = tokens.INTEGER

			return tok
		default:
			tok = tokens.New(tokens.ILLEGAL, string(l.ch))
		}
	}

	l.readChar()

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

func (l *Lexer) skipWhitespaces() {
	for isWhitespace(l.ch) {
		l.readChar()
	}
}

func (l *Lexer) readIdentifier() string {
	position := l.position

	for isLetter(l.ch) {
		l.readChar()
	}

	return l.input[position:l.position]
}

func (l *Lexer) readNumber() string {
	position := l.position

	for isDigit(l.ch) {
		l.readChar()
	}

	return l.input[position:l.position]
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isWhitespace(ch byte) bool {
	return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}
