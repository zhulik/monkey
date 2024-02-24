package ast

import (
	"bytes"

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

func (vn ValueNode[T]) Value() T {
	return vn.V
}

func (vn ValueNode[T]) TokenLiteral() string {
	return vn.Token.Literal()
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
