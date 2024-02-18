package lexer

import (
	"errors"
	"fmt"
	"io"

	"github.com/zhulik/monkey/tokens"
)

var ErrIllegalCharacter = errors.New("illegal character")

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

func (l *Lexer) NextToken() (tokens.Token, error) { //nolint:cyclop,funlen
	if l.position >= len(l.input) {
		return tokens.Token{}, io.EOF
	}

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
		if l.peekChar() == '=' {
			l.readChar()

			tok = tokens.New(tokens.LTE)
		} else {
			tok = tokens.New(tokens.LT)
		}
	case '>':
		if l.peekChar() == '=' {
			l.readChar()

			tok = tokens.New(tokens.GTE)
		} else {
			tok = tokens.New(tokens.GT)
		}
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
		defer l.readChar()

		return tokens.Token{}, io.EOF
	default:
		switch {
		case isLetter(l.ch):
			return l.identifierToken(), nil
		case isDigit(l.ch):
			return l.numberToken(), nil
		default:
			defer l.readChar()

			return tokens.Token{}, fmt.Errorf("%w: '%s'", ErrIllegalCharacter, string(l.ch))
		}
	}

	l.readChar()

	return tok, nil
}

func (l *Lexer) Tokens() ([]tokens.Token, error) {
	tkns := []tokens.Token{}

	for token, err := l.NextToken(); !errors.Is(err, io.EOF); token, err = l.NextToken() {
		if err != nil {
			return []tokens.Token{}, err
		}

		tkns = append(tkns, token)
	}

	return tkns, nil
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
