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

func (l *Lexer) NextToken() tokens.Token { //nolint:cyclop,funlen
	l.skipWhitespaces()

	var tok tokens.Token

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			l.readChar()

			tok = tokens.New(tokens.EQ)
		} else {
			tok = tokens.New(tokens.ASSIGN)
		}
	case '+':
		tok = tokens.New(tokens.PLUS)
	case '-':
		tok = tokens.New(tokens.MINUS)
	case '/':
		tok = tokens.New(tokens.SLASH)
	case '*':
		tok = tokens.New(tokens.ASTERISK)
	case '<':
		tok = tokens.New(tokens.LT)
	case '>':
		tok = tokens.New(tokens.GT)
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
	case '!':
		if l.peekChar() == '=' {
			l.readChar()

			tok = tokens.New(tokens.NEQ)
		} else {
			tok = tokens.New(tokens.BANG)
		}

	case 0:
		tok = tokens.New(tokens.EOF)
	default:
		switch {
		case isLetter(l.ch):
			return l.identifierToken()
		case isDigit(l.ch):
			return l.numberToken()
		default:
			tok = tokens.New(tokens.ILLEGAL, string(l.ch))
		}
	}

	l.readChar()

	return tok
}

func (l *Lexer) identifierToken() tokens.Token {
	literal := l.readIdentifier()

	var token tokens.Token

	token.Type = tokens.TypeFromLiteral(literal)
	if token.Type == tokens.IDENTIFIER {
		token.Literal = literal
	}

	return token
}

func (l *Lexer) numberToken() tokens.Token {
	var token tokens.Token

	token.Literal = l.readNumber()
	token.Type = tokens.INTEGER

	return token
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

func (l Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}

	return l.input[l.readPosition]
}

func (l *Lexer) skipWhitespaces() {
	l.readAll(isWhitespace)
}

func (l *Lexer) readIdentifier() string {
	return l.readAll(isLetter)
}

func (l *Lexer) readNumber() string {
	return l.readAll(isDigit)
}

func (l *Lexer) readAll(fn func(byte) bool) string {
	position := l.position

	for fn(l.ch) {
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
