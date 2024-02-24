package ast

import (
	"bytes"
)

type StatementNode[T any] struct {
	ValueNode[T]
}

func (sn StatementNode[T]) statementNode() {}

type LetStatement struct {
	StatementNode[Expression]
	Name *IdentifierExpression
}

func (l LetStatement) String() string {
	out := bytes.Buffer{}

	out.WriteString(l.TokenLiteral() + " ")
	out.WriteString(l.Name.String() + " = ")

	if l.Value() != nil {
		out.WriteString(l.ValueNode.V.String())
	}

	out.WriteString(";")

	return out.String()
}

type ReturnStatement struct {
	StatementNode[Expression]
}

func (r ReturnStatement) String() string {
	out := bytes.Buffer{}

	out.WriteString(r.TokenLiteral() + " ")

	if r.Value() != nil {
		out.WriteString(r.Value().String())
	}

	out.WriteString(";")

	return out.String()
}

type ExpressionStatement struct {
	StatementNode[Expression]
}

func (e ExpressionStatement) String() string {
	if e.Value() != nil {
		return e.Value().String()
	}

	return ""
}

type BlockStatement struct {
	StatementNode[[]Statement]
}

func (p BlockStatement) String() string {
	var out bytes.Buffer

	for _, stmt := range p.Value() {
		out.WriteString(stmt.String())
	}

	return out.String()
}
