package ast

import "github.com/riadafridishibly/go-monkey/token"

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

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}

	return ""
}

type LetStatement struct {
	Token token.Token // token.LET
	Name  *Identifier
	Value Expression
}

var _ Statement = (*LetStatement)(nil)

func (ls *LetStatement) statementNode()       {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }

type Identifier struct {
	Token token.Token // token.IDENT
	Value string
}

var _ Expression = (*Identifier)(nil)

func (ident *Identifier) expressionNode()      {}
func (ident *Identifier) TokenLiteral() string { return ident.Token.Literal }

type ReturnStatement struct {
	Token       token.Token // token.RETURN
	ReturnValue Expression
}

// TokenLiteral implements Statement.
func (r *ReturnStatement) TokenLiteral() string {
	return r.Token.Literal
}

// statementNode implements Statement.
func (r *ReturnStatement) statementNode() {}

var _ Statement = (*ReturnStatement)(nil)
