package ast

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/samber/lo"
)

type ExpressionNode[T any] struct {
	ValueNode[T]
}

func (sn ExpressionNode[T]) expressionNode() {}

type IdentifierExpression struct {
	ExpressionNode[string]
}

func (i IdentifierExpression) String() string {
	return i.Value()
}

type IntegerExpression struct {
	ExpressionNode[int64]
}

func (i IntegerExpression) String() string {
	return i.Token.Literal()
}

type PrefixExpression struct {
	ExpressionNode[Expression]
	Operator string
}

func (p PrefixExpression) String() string {
	return fmt.Sprintf("(%s%s)", p.Operator, p.Value().String())
}

type InfixExpression struct {
	ExpressionNode[Expression] // Value is the left expression
	Operator                   string
	Right                      Expression
}

func (p InfixExpression) String() string {
	return fmt.Sprintf("(%s %s %s)", p.Value().String(), p.Operator, p.Right.String())
}

type BooleanExpression struct {
	ExpressionNode[bool]
}

func (p BooleanExpression) String() string {
	return p.Token.Literal()
}

type IfExpression struct {
	ExpressionNode[Expression] // Value is the condition expression
	Then                       *BlockStatement
	Else                       *BlockStatement
}

func (p IfExpression) String() string {
	var out bytes.Buffer

	out.WriteString("if ")
	out.WriteString(p.Value().String())

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
	ExpressionNode[*BlockStatement] // Value is the block
	Arguments                       []*IdentifierExpression
}

func (p FunctionExpression) String() string {
	var out bytes.Buffer

	params := lo.Map(p.Arguments, func(item *IdentifierExpression, _ int) string {
		return item.String()
	})

	out.WriteString(p.TokenLiteral() + "(" + strings.Join(params, ", ") + ")")
	body := p.Value().String()

	if len(body) > 0 {
		out.WriteString(" { " + p.Value().String() + " }")
	} else {
		out.WriteString(" { }")
	}

	return out.String()
}

type CallExpression struct {
	ExpressionNode[Expression] // Value is the FunctionExpression or IdentifierExpression
	Arguments                  []Expression
}

func (p CallExpression) String() string {
	var out bytes.Buffer

	args := lo.Map(p.Arguments, func(item Expression, _ int) string {
		return item.String()
	})

	out.WriteString(p.Value().String() + "(")

	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")

	return out.String()
}

type NilExpression struct {
	ExpressionNode[any] // Value is not used
}

func (p NilExpression) String() string {
	return "nil"
}
