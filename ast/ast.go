package ast

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/samber/lo"
	"github.com/zhulik/monkey/tokens"
)

type Node interface {
	TokenLiteral() string
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

func (p Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}

	return ""
}

func (p Program) String() string {
	out := bytes.Buffer{}

	for _, stmt := range p.Statements {
		out.WriteString(stmt.String())
	}

	return out.String()
}

type LetStatement struct {
	Token tokens.Token
	Name  *IdentifierExpression
	Value Expression
}

func (l LetStatement) statementNode() {}
func (l LetStatement) TokenLiteral() string {
	return l.Token.Literal()
}

func (l LetStatement) String() string {
	out := bytes.Buffer{}

	out.WriteString(l.TokenLiteral() + " ")
	out.WriteString(l.Name.String() + " = ")

	if l.Value != nil {
		out.WriteString(l.Value.String())
	}

	out.WriteString(";")

	return out.String()
}

type IdentifierExpression struct {
	Token tokens.Token
	Value string
}

func (i IdentifierExpression) expressionNode() {}
func (i IdentifierExpression) TokenLiteral() string {
	return i.Token.Literal()
}

func (i IdentifierExpression) String() string {
	return i.Value
}

type ReturnStatement struct {
	Token tokens.Token
	Value Expression
}

func (r ReturnStatement) statementNode() {}
func (r ReturnStatement) TokenLiteral() string {
	return r.Token.Literal()
}

func (r ReturnStatement) String() string {
	out := bytes.Buffer{}

	out.WriteString(r.TokenLiteral() + " ")

	if r.Value != nil {
		out.WriteString(r.Value.String())
	}

	out.WriteString(";")

	return out.String()
}

type ExpressionStatement struct {
	Token      tokens.Token
	Expression Expression
}

func (e ExpressionStatement) statementNode() {}
func (e ExpressionStatement) TokenLiteral() string {
	return e.Token.Literal()
}

func (e ExpressionStatement) String() string {
	if e.Expression != nil {
		return e.Expression.String()
	}

	return ""
}

type IntegerExpression struct {
	Token tokens.Token
	Value int64
}

func (i IntegerExpression) expressionNode() {}
func (i IntegerExpression) TokenLiteral() string {
	return i.Token.Literal()
}

func (i IntegerExpression) String() string {
	return i.Token.Literal()
}

type PrefixExpression struct {
	Token    tokens.Token
	Operator string
	Right    Expression
}

func (p PrefixExpression) expressionNode() {}
func (p PrefixExpression) TokenLiteral() string {
	return p.Token.Literal()
}

func (p PrefixExpression) String() string {
	return fmt.Sprintf("(%s%s)", p.Operator, p.Right.String())
}

type InfixExpression struct {
	Token    tokens.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (p InfixExpression) expressionNode() {}
func (p InfixExpression) TokenLiteral() string {
	return p.Token.Literal()
}

func (p InfixExpression) String() string {
	return fmt.Sprintf("(%s %s %s)", p.Left.String(), p.Operator, p.Right.String())
}

type BooleanExpression struct {
	Token tokens.Token
	Value bool
}

func (p BooleanExpression) expressionNode() {}
func (p BooleanExpression) TokenLiteral() string {
	return p.Token.Literal()
}

func (p BooleanExpression) String() string {
	return p.Token.Literal()
}

type IfExpression struct {
	Token     tokens.Token
	Condition Expression
	Then      *BlockStatement
	Else      *BlockStatement
}

func (p IfExpression) expressionNode() {}
func (p IfExpression) TokenLiteral() string {
	return p.Token.Literal()
}

func (p IfExpression) String() string {
	var out bytes.Buffer

	out.WriteString("if ")
	out.WriteString(p.Condition.String())

	thenBody := p.Then.String()

	if len(thenBody) > 0 {
		out.WriteString(" { " + thenBody + " }")
	} else {
		out.WriteString(" { }")
	}

	if p.Else != nil {
		elseBody := p.Else.String()

		if len(elseBody) > 0 {
			out.WriteString(" else { " + elseBody + " }")
		} else {
			out.WriteString(" else { }")
		}
	}

	return out.String()
}

type BlockStatement struct {
	Token      tokens.Token
	Statements []Statement
}

func (p BlockStatement) statementNode() {}
func (p BlockStatement) TokenLiteral() string {
	return p.Token.Literal()
}

func (p BlockStatement) String() string {
	var out bytes.Buffer

	for _, stmt := range p.Statements {
		out.WriteString(stmt.String())
	}

	return out.String()
}

type FunctionExpression struct {
	Token      tokens.Token
	Parameters []*IdentifierExpression
	Body       *BlockStatement
}

func (p FunctionExpression) expressionNode() {}
func (p FunctionExpression) TokenLiteral() string {
	return p.Token.Literal()
}

func (p FunctionExpression) String() string {
	var out bytes.Buffer

	params := lo.Map(p.Parameters, func(item *IdentifierExpression, _ int) string {
		return item.String()
	})

	out.WriteString(p.TokenLiteral() + "(" + strings.Join(params, ", ") + ")")
	body := p.Body.String()

	if len(body) > 0 {
		out.WriteString(" { " + p.Body.String() + " }")
	} else {
		out.WriteString(" { }")
	}

	return out.String()
}
