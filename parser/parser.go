package parser

import (
	"errors"
	"fmt"
	"io"
	"strconv"

	"github.com/zhulik/monkey/ast"
	"github.com/zhulik/monkey/lexer"
	"github.com/zhulik/monkey/tokens"
)

var (
	ErrInvalidToken        = errors.New("invalid token")
	ErrNoPrefixParserFound = errors.New("no prefix parse function  found for")

	precedences = map[tokens.TokenType]int{ //nolint:gochecknoglobals
		tokens.EQ:       EQUALS,
		tokens.NEQ:      EQUALS,
		tokens.LT:       LESSGREATER,
		tokens.GT:       LESSGREATER,
		tokens.LTE:      LESSGREATER, // TODO: check precedences
		tokens.GTE:      LESSGREATER, // TODO: check precedences
		tokens.PLUS:     SUM,
		tokens.MINUS:    SUM,
		tokens.SLASH:    PRODUCT,
		tokens.ASTERISK: PRODUCT,
	}
)

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

const (
	_ int = iota

	LOWEST
	EQUALS
	LESSGREATER
	SUM
	PRODUCT
	PREFIX
	CALL
)

type Parser struct {
	lexer        *lexer.Lexer
	currentToken tokens.Token
	peekToken    tokens.Token

	prefixParseFns map[tokens.TokenType]prefixParseFn
	infixParseFns  map[tokens.TokenType]infixParseFn
}

func New(l *lexer.Lexer) (*Parser, error) {
	parser := &Parser{lexer: l}

	parser.prefixParseFns = map[tokens.TokenType]prefixParseFn{
		tokens.IDENTIFIER: parser.parseIdentifierExpression,
		tokens.INTEGER:    parser.parseIntegerExpression,
		tokens.BANG:       parser.parsePrefixExpression,
		tokens.MINUS:      parser.parsePrefixExpression,
	}
	parser.infixParseFns = map[tokens.TokenType]infixParseFn{
		tokens.PLUS:     parser.parseInfixExpression,
		tokens.MINUS:    parser.parseInfixExpression,
		tokens.SLASH:    parser.parseInfixExpression,
		tokens.ASTERISK: parser.parseInfixExpression,
		tokens.EQ:       parser.parseInfixExpression,
		tokens.NEQ:      parser.parseInfixExpression,
		tokens.LT:       parser.parseInfixExpression,
		tokens.GT:       parser.parseInfixExpression,
		tokens.LTE:      parser.parseInfixExpression,
		tokens.GTE:      parser.parseInfixExpression,
	}

	err := parser.nextToken()
	if err != nil {
		return nil, err
	}

	err = parser.nextToken()
	if err != nil {
		return nil, err
	}

	return parser, nil
}

func (p *Parser) ParseProgram() (*ast.Program, error) {
	program := &ast.Program{}

	for {
		stmt, err := p.parseStatement()
		if err != nil {
			return nil, err
		}

		program.Statements = append(program.Statements, stmt)

		err = p.nextToken()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			return nil, err
		}
	}

	return program, nil
}

func (p *Parser) parseStatement() (ast.Statement, error) {
	switch p.currentToken.Type { //nolint:exhaustive
	case tokens.LET:
		return p.parseLetStatement()
	case tokens.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseLetStatement() (*ast.LetStatement, error) {
	stmt := &ast.LetStatement{Token: p.currentToken}

	err := p.expectPeek(tokens.IDENTIFIER)
	if err != nil {
		return nil, err
	}

	stmt.Name = &ast.IdentifierExpression{Token: p.currentToken, Value: p.currentToken.Literal()}

	err = p.expectPeek(tokens.ASSIGN)
	if err != nil {
		return nil, err
	}

	for p.currentToken.Type != tokens.SEMICOLON {
		nErr := p.nextToken()
		if nErr != nil {
			if errors.Is(nErr, io.EOF) {
				return stmt, nil
			}

			return nil, nErr
		}
	}

	return stmt, nil
}

func (p *Parser) parseReturnStatement() (*ast.ReturnStatement, error) {
	stmt := &ast.ReturnStatement{Token: p.currentToken}

	err := p.nextToken()
	if err != nil {
		return nil, err
	}

	for p.currentToken.Type != tokens.SEMICOLON {
		nErr := p.nextToken()
		if nErr != nil {
			if errors.Is(nErr, io.EOF) {
				return stmt, nil
			}

			return nil, nErr
		}
	}

	return stmt, nil
}

func (p *Parser) parseExpressionStatement() (*ast.ExpressionStatement, error) {
	stmt := &ast.ExpressionStatement{Token: p.currentToken}

	expression, err := p.parseExpression(LOWEST)
	if err != nil {
		return nil, err
	}

	stmt.Expression = expression

	if p.peekToken.Type == tokens.SEMICOLON {
		nErr := p.nextToken()
		if nErr != nil {
			if errors.Is(nErr, io.EOF) {
				return stmt, nil
			}

			return nil, nErr
		}
	}

	return stmt, nil
}

func (p *Parser) parseExpression(precedence int) (ast.Expression, error) {
	prefix := p.prefixParseFns[p.currentToken.Type]
	if prefix == nil {
		return nil, fmt.Errorf("%w %s", ErrNoPrefixParserFound, p.currentToken.Type)
	}

	leftExpr := prefix()

	for p.currentToken.Type != tokens.SEMICOLON && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExpr, nil
		}

		nErr := p.nextToken()
		if nErr != nil {
			if errors.Is(nErr, io.EOF) {
				return leftExpr, nil
			}

			return nil, nErr
		}

		leftExpr = infix(leftExpr)
	}

	return leftExpr, nil
}

func (p *Parser) parseIdentifierExpression() ast.Expression {
	return &ast.IdentifierExpression{Token: p.currentToken, Value: p.currentToken.Literal()}
}

func (p *Parser) parseIntegerExpression() ast.Expression {
	expr := &ast.IntegerExpression{Token: p.currentToken}

	value, err := strconv.ParseInt(p.currentToken.Literal(), 10, 64)
	if err != nil {
		panic(err) // TODO: fix me
	}

	expr.Value = value

	return expr
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expr := &ast.PrefixExpression{Token: p.currentToken, Operator: p.currentToken.Literal()}

	err := p.nextToken()
	if err != nil {
		panic(err) // TODO: fix me
	}

	right, err := p.parseExpression(PREFIX)
	if err != nil {
		panic(err) // TODO: fix me
	}

	expr.Right = right

	return expr
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expr := &ast.InfixExpression{
		Token:    p.currentToken,
		Operator: p.currentToken.Literal(),
		Left:     left,
	}

	precedence := p.currentPrecedence()

	err := p.nextToken()
	if err != nil {
		panic(err) // TODO: fix me
	}

	right, err := p.parseExpression(precedence)
	if err != nil {
		panic(err) // TODO: fix me
	}

	expr.Right = right

	return expr
}

func (p *Parser) expectPeek(tokenType tokens.TokenType) error {
	if p.peekToken.Type == tokenType {
		err := p.nextToken()
		if err != nil {
			return err
		}

		return nil
	}

	return fmt.Errorf("%w. Expected: %s, found: %s", ErrInvalidToken, tokenType, p.peekToken.Type)
}

func (p *Parser) nextToken() error {
	p.currentToken = p.peekToken

	peekToken, err := p.lexer.NextToken()
	if err != nil {
		return fmt.Errorf("reading next token error: %w", err)
	}

	p.peekToken = peekToken

	return nil
}

func (p *Parser) peekPrecedence() int {
	return precedence(p.peekToken)
}

func (p *Parser) currentPrecedence() int {
	return precedence(p.currentToken)
}

func precedence(token tokens.Token) int {
	if p, ok := precedences[token.Type]; ok {
		return p
	}

	return LOWEST
}
