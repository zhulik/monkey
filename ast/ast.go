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

type ValueNode[T any] struct {
	tokens.Token
	V T
}

func (v ValueNode[T]) Value() T {
	return v.V
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

type IdentifierExpression struct {
	tokens.Token
	Value string
}

func (i IdentifierExpression) expressionNode() {}
func (i IdentifierExpression) TokenLiteral() string {
	return i.Token.Literal()
}

func (i IdentifierExpression) String() string {
	return i.Value
}

type IntegerExpression struct {
	tokens.Token
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
	tokens.Token
	Operator string
	Value    Expression
}

func (p PrefixExpression) expressionNode() {}
func (p PrefixExpression) TokenLiteral() string {
	return p.Token.Literal()
}

func (p PrefixExpression) String() string {
	return fmt.Sprintf("(%s%s)", p.Operator, p.Value.String())
}

type InfixExpression struct {
	tokens.Token
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
	tokens.Token
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
	tokens.Token
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

type FunctionExpression struct {
	tokens.Token
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

type CallExpression struct {
	tokens.Token
	Function  Expression // FunctionExpression or IdentifierExpression
	Arguments []Expression
}

func (p CallExpression) expressionNode() {}
func (p CallExpression) TokenLiteral() string {
	return p.Token.Literal()
}

func (p CallExpression) String() string {
	var out bytes.Buffer

	args := lo.Map(p.Arguments, func(item Expression, _ int) string {
		return item.String()
	})

	out.WriteString(p.Function.String() + "(")

	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")

	return out.String()
}
