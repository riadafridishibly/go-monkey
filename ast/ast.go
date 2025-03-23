package ast

import (
	"fmt"
	"strings"

	"github.com/riadafridishibly/go-monkey/token"
)

type Node interface {
	TokenLiteral() string
	fmt.Stringer
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

// String implements Node.
func (p *Program) String() string {
	sb := strings.Builder{}
	for _, s := range p.Statements {
		sb.WriteString(s.String())
	}
	return sb.String()
}

var _ Node = (*Program)(nil)

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

// String implements Statement.
func (ls *LetStatement) String() string {
	sb := strings.Builder{}

	sb.WriteString(ls.TokenLiteral() + " ")
	sb.WriteString(ls.Name.String())
	sb.WriteString(" = ")

	if ls.Value != nil {
		sb.WriteString(ls.Value.String())
	}

	sb.WriteString(";")

	return sb.String()
}

var _ Statement = (*LetStatement)(nil)

func (ls *LetStatement) statementNode()       {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }

type Identifier struct {
	Token token.Token // token.IDENT
	Value string
}

// String implements Expression.
func (ident *Identifier) String() string {
	return ident.Value
}

var _ Expression = (*Identifier)(nil)

func (ident *Identifier) expressionNode()      {}
func (ident *Identifier) TokenLiteral() string { return ident.Token.Literal }

type ReturnStatement struct {
	Token       token.Token // token.RETURN
	ReturnValue Expression
}

// String implements Statement.
func (r *ReturnStatement) String() string {
	sb := strings.Builder{}

	sb.WriteString(r.TokenLiteral() + " ")
	if r.ReturnValue != nil {
		sb.WriteString(r.ReturnValue.String())
	}
	sb.WriteString(";")

	return sb.String()
}

// TokenLiteral implements Statement.
func (r *ReturnStatement) TokenLiteral() string {
	return r.Token.Literal
}

// statementNode implements Statement.
func (r *ReturnStatement) statementNode() {}

var _ Statement = (*ReturnStatement)(nil)

type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

// String implements Statement.
func (e *ExpressionStatement) String() string {
	if e.Expression != nil {
		return e.Expression.String()
	}
	return ""
}

// TokenLiteral implements Statement.
func (e *ExpressionStatement) TokenLiteral() string {
	return e.Token.Literal
}

// statementNode implements Statement.
func (e *ExpressionStatement) statementNode() {}

var _ Statement = (*ExpressionStatement)(nil)

type InetegerLiteral struct {
	Token token.Token
	Value int64
}

// String implements Expression.
func (i *InetegerLiteral) String() string {
	return i.Token.Literal
}

// TokenLiteral implements Expression.
func (i *InetegerLiteral) TokenLiteral() string {
	return i.Token.Literal
}

// expressionNode implements Expression.
func (i *InetegerLiteral) expressionNode() {}

var _ Expression = (*InetegerLiteral)(nil)

type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

// String implements Expression.
func (p *PrefixExpression) String() string {
	sb := strings.Builder{}
	sb.WriteString("(")
	sb.WriteString(p.Operator)
	sb.WriteString(p.Right.String())
	sb.WriteString(")")
	return sb.String()
}

// TokenLiteral implements Expression.
func (p *PrefixExpression) TokenLiteral() string {
	return p.Token.Literal
}

// expressionNode implements Expression.
func (p *PrefixExpression) expressionNode() {}

var _ Expression = (*PrefixExpression)(nil)

type InfixExpression struct {
	Token    token.Token
	Left     Expression
	Operator string
	Right    Expression
}

// String implements Expression.
func (i *InfixExpression) String() string {
	sb := strings.Builder{}
	sb.WriteString("(")
	sb.WriteString(i.Left.String())
	sb.WriteString(" " + i.Operator + " ")
	sb.WriteString(i.Right.String())
	sb.WriteString(")")
	return sb.String()
}

// TokenLiteral implements Expression.
func (i *InfixExpression) TokenLiteral() string {
	return i.Token.Literal
}

// expressionNode implements Expression.
func (i *InfixExpression) expressionNode() {}

var _ Expression = (*InfixExpression)(nil)
