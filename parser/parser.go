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

func New(l *lexer.Lexer) *Parser {
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

	return parser
}

func (p *Parser) ParseProgram() (*ast.Program, error) {
	program := &ast.Program{}

	nErr := p.nextToken()
	if nErr != nil {
		return nil, nErr
	}

	for err := p.nextTokenIgnoreEOF(); !errors.Is(err, io.EOF); err = p.nextToken() {
		if err != nil {
			return nil, err
		}

		stmt, sErr := p.parseStatement()
		if sErr != nil {
			return nil, sErr
		}

		program.Statements = append(program.Statements, stmt)
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
	stmt := ast.NewValueNode[ast.LetStatement, ast.Expression](p.currentToken)

	err := p.expectPeek(tokens.IDENTIFIER)
	if err != nil {
		return nil, err
	}

	stmt.Name = ast.NewValueNode[ast.IdentifierExpression](p.currentToken, p.currentToken.Literal())

	err = p.expectPeek(tokens.ASSIGN)
	if err != nil {
		return nil, err
	}

	err = p.nextToken()
	if err != nil {
		return nil, err
	}

	stmt.V, err = p.parseExpression(LOWEST)
	if err != nil {
		return nil, err
	}

	if p.peekToken.Type == tokens.SEMICOLON {
		err = p.nextTokenIgnoreEOF()
		if err != nil {
			return nil, err
		}
	}

	return stmt, nil
}

func (p *Parser) parseReturnStatement() (*ast.ReturnStatement, error) {
	stmt := ast.NewValueNode[ast.ReturnStatement, ast.Expression](p.currentToken)

	err := p.nextToken()
	if err != nil {
		return nil, err
	}

	stmt.V, err = p.parseExpression(LOWEST)
	if err != nil {
		return nil, err
	}

	if p.peekToken.Type == tokens.SEMICOLON {
		err = p.nextTokenIgnoreEOF()
		if err != nil {
			return nil, err
		}
	}

	return stmt, nil
}

func (p *Parser) parseExpressionStatement() (*ast.ExpressionStatement, error) {
	stmt := ast.NewValueNode[ast.ExpressionStatement, ast.Expression](p.currentToken)

	var err error

	stmt.V, err = p.parseExpression(LOWEST)
	if err != nil {
		return nil, err
	}

	if p.peekToken.Type == tokens.SEMICOLON {
		nErr := p.nextTokenIgnoreEOF()
		if nErr != nil {
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

		nErr := p.nextTokenIgnoreEOF()
		if nErr != nil {
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
	return ast.NewValueNode[ast.IdentifierExpression](p.currentToken, p.currentToken.Literal()), nil
}

func (p *Parser) parseIntegerExpression() (ast.Expression, error) {
	expr := ast.NewValueNode[ast.IntegerExpression, int64](p.currentToken)

	var err error

	expr.V, err = strconv.ParseInt(p.currentToken.Literal(), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("error parsing integer expression: %w", err)
	}

	return expr, nil
}

func (p *Parser) parsePrefixExpression() (ast.Expression, error) {
	expr := ast.NewValueNode[ast.PrefixExpression, ast.Expression](p.currentToken)
	expr.Operator = p.currentToken.Literal()

	err := p.nextTokenIgnoreEOF()
	if err != nil {
		return nil, err
	}

	expr.V, err = p.parseExpression(PREFIX)
	if err != nil {
		return nil, err
	}

	return expr, nil
}

func (p *Parser) parseBooleanExpression() (ast.Expression, error) {
	return ast.NewValueNode[ast.BooleanExpression](p.currentToken, p.currentToken.Type == tokens.TRUE), nil
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

	err = ignoreError(p.expectPeek(tokens.RPAREN), io.EOF)
	if err != nil {
		return nil, err
	}

	return expr, nil
}

func (p *Parser) parseIfExpression() (ast.Expression, error) { //nolint:cyclop
	expr := ast.NewValueNode[ast.IfExpression, ast.Expression](p.currentToken)

	err := p.expectPeek(tokens.LPAREN)
	if err != nil {
		return nil, err
	}

	err = p.nextToken()
	if err != nil {
		return nil, err
	}

	expr.V, err = p.parseExpression(LOWEST)
	if err != nil {
		return nil, err
	}

	err = p.expectPeek(tokens.RPAREN)
	if err != nil {
		return nil, err
	}

	err = p.expectPeek(tokens.LBRACE)
	if err != nil {
		return nil, err
	}

	expr.Then, err = p.parseBlockStatement()
	if err != nil {
		return nil, err
	}

	if p.peekToken.Type == tokens.ELSE {
		err = p.nextToken()
		if err != nil {
			return nil, err
		}

		err = p.expectPeek(tokens.LBRACE)
		if err != nil {
			return nil, err
		}

		expr.Else, err = p.parseBlockStatement()
		if err != nil {
			return nil, err
		}
	}

	return expr, nil
}

func (p *Parser) parseBlockStatement() (*ast.BlockStatement, error) {
	block := ast.NewValueNode[ast.BlockStatement, []ast.Statement](p.currentToken)

	err := p.nextTokenIgnoreEOF()
	if err != nil {
		return nil, err
	}

	for p.currentToken.Type != tokens.RBRACE {
		stmt, sErr := p.parseStatement()
		if sErr != nil {
			return nil, sErr
		}

		block.V = append(block.V, stmt)

		nErr := p.nextTokenIgnoreEOF()
		if nErr != nil {
			return nil, nErr
		}
	}

	return block, nil
}

func (p *Parser) parseInfixExpression(left ast.Expression) (ast.Expression, error) {
	expr := ast.NewValueNode[ast.InfixExpression](p.currentToken, left)
	expr.Operator = p.currentToken.Literal()

	precedence := precedence(p.currentToken)

	err := p.nextTokenIgnoreEOF()
	if err != nil {
		return nil, err
	}

	expr.Right, err = p.parseExpression(precedence)
	if err != nil {
		return nil, err
	}

	return expr, nil
}

func (p *Parser) parseCallExpression(function ast.Expression) (ast.Expression, error) {
	expr := ast.NewValueNode[ast.CallExpression](p.currentToken, function)

	var err error

	expr.Arguments, err = p.parseCallArguments()
	if err != nil {
		return nil, err
	}

	return expr, nil
}

func (p *Parser) parseCallArguments() ([]ast.Expression, error) {
	args := []ast.Expression{}

	for err := p.nextToken(); p.currentToken.Type != tokens.RPAREN; err = p.nextToken() {
		if err != nil {
			return nil, err
		}

		if p.currentToken.Type != tokens.COMMA {
			expr, eErr := p.parseExpression(LOWEST)
			if eErr != nil {
				return nil, eErr
			}

			args = append(args, expr)
		}
	}

	return args, nil
}

func (p *Parser) parseFunctionExpression() (ast.Expression, error) {
	expr := ast.NewValueNode[ast.FunctionExpression, *ast.BlockStatement](p.currentToken)

	err := p.expectPeek(tokens.LPAREN)
	if err != nil {
		return nil, err
	}

	params, err := p.parseFunctionArguments()
	if err != nil {
		return nil, err
	}

	expr.Parameters = params

	err = p.expectPeek(tokens.LBRACE)
	if err != nil {
		return nil, err
	}

	expr.V, err = p.parseBlockStatement()
	if err != nil {
		return nil, err
	}

	return expr, nil
}

func (p *Parser) parseFunctionArguments() ([]*ast.IdentifierExpression, error) {
	args := []*ast.IdentifierExpression{}

	for err := p.nextToken(); p.currentToken.Type != tokens.RPAREN; err = p.nextToken() {
		if err != nil {
			return nil, err
		}

		if p.currentToken.Type != tokens.COMMA {
			args = append(args, ast.NewValueNode[ast.IdentifierExpression](p.currentToken, p.currentToken.Literal()))
		}
	}

	return args, nil
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

func (p *Parser) nextTokenIgnoreEOF() error {
	return ignoreError(p.nextToken(), io.EOF)
}

func ignoreError(err, errToIgnore error) error {
	if errors.Is(err, errToIgnore) {
		return nil
	}

	return err
}

func precedence(token tokens.Token) int {
	if p, ok := precedences[token.Type]; ok {
		return p
	}

	return LOWEST
}
