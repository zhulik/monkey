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
	ErrNoPrefixParserFound = errors.New("no prefix parse function found for")

	precedences = map[tokens.TokenType]int{ //nolint:gochecknoglobals
		tokens.EQ:       EQUALS,
		tokens.NEQ:      EQUALS,
		tokens.LT:       LESSGREATER,
		tokens.GT:       LESSGREATER,
		tokens.LTE:      LESSGREATER,
		tokens.GTE:      LESSGREATER,
		tokens.PLUS:     SUM,
		tokens.MINUS:    SUM,
		tokens.SLASH:    PRODUCT,
		tokens.ASTERISK: PRODUCT,
		tokens.LPAREN:   CALL,
	}
)

type (
	prefixParseFn func() (ast.Expression, error)
	infixParseFn  func(ast.Expression) (ast.Expression, error)
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
		tokens.TRUE:       parser.parseBooleanExpression,
		tokens.FALSE:      parser.parseBooleanExpression,
		tokens.LPAREN:     parser.parseGroupedExpression,
		tokens.IF:         parser.parseIfExpression,
		tokens.FUNCTION:   parser.parseFunctionExpression,
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
		tokens.LPAREN:   parser.parseCallExpression,
	}

	err := parser.nextToken()
	if err != nil {
		return nil, err
	}

	// It's possible the have only one token in the program.
	err = parser.nextToken()
	if err != nil {
		if !errors.Is(err, io.EOF) {
			return nil, err
		}
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

	err = p.nextToken()
	if err != nil {
		return nil, err
	}

	expr, err := p.parseExpression(LOWEST)
	if err != nil {
		return nil, err
	}

	stmt.Value = expr

	if p.peekToken.Type == tokens.SEMICOLON {
		err = p.nextToken()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return stmt, nil
			}

			return nil, err
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

	value, err := p.parseExpression(LOWEST)
	if err != nil {
		return nil, err
	}

	stmt.Value = value

	if p.peekToken.Type == tokens.SEMICOLON {
		err = p.nextToken()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return stmt, nil
			}

			return nil, err
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

func (p *Parser) parseExpression(prec int) (ast.Expression, error) {
	prefix := p.prefixParseFns[p.currentToken.Type]
	if prefix == nil {
		return nil, fmt.Errorf("%w %s", ErrNoPrefixParserFound, p.currentToken.Type)
	}

	leftExpr, err := prefix()
	if err != nil {
		return nil, err
	}

	for p.currentToken.Type != tokens.SEMICOLON && prec < precedence(p.peekToken) {
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

		var lErr error

		leftExpr, lErr = infix(leftExpr)
		if lErr != nil {
			if errors.Is(lErr, io.EOF) {
				return leftExpr, nil
			}

			return nil, lErr
		}
	}

	return leftExpr, nil
}

func (p *Parser) parseIdentifierExpression() (ast.Expression, error) {
	return &ast.IdentifierExpression{Token: p.currentToken, Value: p.currentToken.Literal()}, nil
}

func (p *Parser) parseIntegerExpression() (ast.Expression, error) {
	expr := &ast.IntegerExpression{Token: p.currentToken}

	value, err := strconv.ParseInt(p.currentToken.Literal(), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("error parsing integer expression: %w", err)
	}

	expr.Value = value

	return expr, nil
}

func (p *Parser) parsePrefixExpression() (ast.Expression, error) {
	expr := &ast.PrefixExpression{Token: p.currentToken, Operator: p.currentToken.Literal()}

	err := p.nextToken()
	if err != nil {
		if !errors.Is(err, io.EOF) {
			return nil, err
		}
	}

	right, err := p.parseExpression(PREFIX)
	if err != nil {
		return nil, err
	}

	expr.Right = right

	return expr, nil
}

func (p *Parser) parseBooleanExpression() (ast.Expression, error) {
	return ast.BooleanExpression{Token: p.currentToken, Value: p.currentToken.Type == tokens.TRUE}, nil
}

func (p *Parser) parseGroupedExpression() (ast.Expression, error) {
	err := p.nextToken()
	if err != nil {
		return nil, err
	}

	expr, err := p.parseExpression(LOWEST)
	if err != nil {
		return nil, err
	}

	err = p.expectPeek(tokens.RPAREN)
	if err != nil {
		if errors.Is(err, io.EOF) {
			return expr, nil
		}

		return nil, err
	}

	return expr, nil
}

func (p *Parser) parseIfExpression() (ast.Expression, error) { //nolint:cyclop
	expr := ast.IfExpression{Token: p.currentToken}

	err := p.expectPeek(tokens.LPAREN)
	if err != nil {
		return nil, err
	}

	err = p.nextToken()
	if err != nil {
		return nil, err
	}

	cond, err := p.parseExpression(LOWEST)
	if err != nil {
		return nil, err
	}

	expr.Condition = cond

	err = p.expectPeek(tokens.RPAREN)
	if err != nil {
		return nil, err
	}

	err = p.expectPeek(tokens.LBRACE)
	if err != nil {
		return nil, err
	}

	block, err := p.parseBlockStatement()
	if err != nil {
		return nil, err
	}

	expr.Then = block

	if p.peekToken.Type == tokens.ELSE {
		err = p.nextToken()
		if err != nil {
			return nil, err
		}

		err = p.expectPeek(tokens.LBRACE)
		if err != nil {
			return nil, err
		}

		block, err = p.parseBlockStatement()
		if err != nil {
			return nil, err
		}

		expr.Else = block
	}

	return expr, nil
}

func (p *Parser) parseBlockStatement() (*ast.BlockStatement, error) {
	block := &ast.BlockStatement{Token: p.currentToken}

	err := p.nextToken()
	if err != nil {
		if errors.Is(err, io.EOF) {
			return block, nil
		}

		return nil, err
	}

	for p.currentToken.Type != tokens.RBRACE {
		stmt, sErr := p.parseStatement()
		if sErr != nil {
			return nil, sErr
		}

		block.Statements = append(block.Statements, stmt)

		nErr := p.nextToken()
		if nErr != nil {
			if errors.Is(nErr, io.EOF) {
				return block, nil
			}

			return nil, nErr
		}
	}

	return block, nil
}

func (p *Parser) parseInfixExpression(left ast.Expression) (ast.Expression, error) {
	expr := &ast.InfixExpression{
		Token:    p.currentToken,
		Operator: p.currentToken.Literal(),
		Left:     left,
	}

	precedence := precedence(p.currentToken)

	err := p.nextToken()
	if err != nil {
		if !errors.Is(err, io.EOF) {
			return nil, err
		}
	}

	right, err := p.parseExpression(precedence)
	if err != nil {
		return nil, err
	}

	expr.Right = right

	return expr, nil
}

func (p *Parser) parseCallExpression(function ast.Expression) (ast.Expression, error) {
	expr := &ast.CallExpression{Token: p.currentToken, Function: function}

	args, err := p.parseCallArguments()
	if err != nil {
		return nil, err
	}

	expr.Arguments = args

	return expr, nil
}

func (p *Parser) parseCallArguments() ([]ast.Expression, error) { //nolint:cyclop
	args := []ast.Expression{}

	if p.peekToken.Type == tokens.RPAREN {
		err := p.nextToken()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return args, nil
			}

			return nil, err
		}

		return args, nil
	}

	err := p.nextToken()
	if err != nil {
		return nil, err
	}

	expr, err := p.parseExpression(LOWEST)
	if err != nil {
		return nil, err
	}

	args = append(args, expr)

	for p.peekToken.Type == tokens.COMMA {
		err = p.nextToken()
		if err != nil {
			return nil, err
		}

		err = p.nextToken()
		if err != nil {
			return nil, err
		}

		expr, err = p.parseExpression(LOWEST)
		if err != nil {
			return nil, err
		}

		args = append(args, expr)
	}

	err = p.expectPeek(tokens.RPAREN)
	if err != nil {
		if errors.Is(err, io.EOF) {
			return args, nil
		}

		return nil, err
	}

	return args, nil
}

func (p *Parser) parseFunctionExpression() (ast.Expression, error) {
	expr := ast.FunctionExpression{Token: p.currentToken}

	err := p.expectPeek(tokens.LPAREN)
	if err != nil {
		return nil, err
	}

	params, err := p.parseFunctionParameters()
	if err != nil {
		return nil, err
	}

	expr.Parameters = params

	err = p.expectPeek(tokens.LBRACE)
	if err != nil {
		return nil, err
	}

	block, err := p.parseBlockStatement()
	if err != nil {
		return nil, err
	}

	expr.Body = block

	return expr, nil
}

func (p *Parser) parseFunctionParameters() ([]*ast.IdentifierExpression, error) {
	params := []*ast.IdentifierExpression{}

	if p.peekToken.Type == tokens.RPAREN {
		err := p.nextToken()
		if err != nil {
			return nil, err
		}

		return params, nil
	}

	err := p.nextToken()
	if err != nil {
		return nil, err
	}

	identifier := &ast.IdentifierExpression{Token: p.currentToken, Value: p.currentToken.Literal()}
	params = append(params, identifier)

	for p.peekToken.Type == tokens.COMMA {
		err = p.nextToken()
		if err != nil {
			return nil, err
		}

		err = p.nextToken()
		if err != nil {
			return nil, err
		}

		identifier = &ast.IdentifierExpression{Token: p.currentToken, Value: p.currentToken.Literal()}
		params = append(params, identifier)
	}

	err = p.expectPeek(tokens.RPAREN)
	if err != nil {
		return nil, err
	}

	return params, nil
}

func (p *Parser) expectPeek(tokenType tokens.TokenType) error {
	if p.peekToken.Type == tokenType {
		err := p.nextToken()
		if err != nil {
			return err
		}

		return nil
	}

	return fmt.Errorf("%w. Expected: %s, found: %s(%s)",
		ErrInvalidToken,
		tokenType,
		p.peekToken.Type,
		p.peekToken.Literal(),
	)
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

func precedence(token tokens.Token) int {
	if p, ok := precedences[token.Type]; ok {
		return p
	}

	return LOWEST
}
