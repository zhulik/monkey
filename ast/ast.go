package ast

import "github.com/zhulik/monkey/tokens"

type Node interface {
	TokenLiteral() string
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

type LetStatement struct {
	Token tokens.Token
	Name  *Identifier
	Value Expression
}

func (l LetStatement) statementNode() {}
func (l LetStatement) TokenLiteral() string {
	return l.Token.Literal
}

type Identifier struct {
	Token tokens.Token
	Value string
}

func (i Identifier) expressionNode() {}
func (i Identifier) TokenLiteral() string {
	return i.Token.Literal
}
