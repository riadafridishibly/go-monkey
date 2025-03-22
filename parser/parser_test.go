package parser

import (
	"testing"

	"github.com/riadafridishibly/go-monkey/ast"
	"github.com/riadafridishibly/go-monkey/lexer"
	"github.com/stretchr/testify/require"
)

func TestLetStatements(t *testing.T) {
	req := require.New(t)
	input := `
let x = 5;
let y = 10;
let foobar = 838383;
`
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	req.NotNil(program, "Program value is nil")
	req.Empty(p.Errors(), "Expected no errors")
	req.Len(program.Statements, 3, "program.Statements should contain 3 statements")

	tests := []struct {
		expectedItent string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		testLetStatement(req, stmt, tt.expectedItent)
	}
}

func testLetStatement(req *require.Assertions, s ast.Statement, name string) {
	req.Equal(s.TokenLiteral(), "let", "Expected TokenLiteral to be let")
	letStmt, ok := s.(*ast.LetStatement)
	req.True(ok, "s not a *ast.LetStatement")
	req.Equal(letStmt.Name.Value, name, "letStmt.Name.Value should be equal")
	req.Equal(letStmt.Name.TokenLiteral(), name)
}

func TestReturnStatement(t *testing.T) {
	input := `
return 5;
return 10;
return 9999;
`
	l := lexer.New(input)
	p := New(l)

	prog := p.ParseProgram()
	if prog == nil {
		t.Fatal("nil program!")
	}
	checkErrors(t, p)

	if len(prog.Statements) != 3 {
		t.Fatalf("expected 3 statements got %v", len(prog.Statements))
	}

	for _, stmt := range prog.Statements {
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("expected stmt *ast.ReturnStatement, but got %T", stmt)
			continue
		}
		if returnStmt.TokenLiteral() != "return" {
			t.Errorf("expected returnStmt.TokenLiteral() = return but got %v", returnStmt.TokenLiteral())
		}
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"

	l := lexer.New(input)
	p := New(l)
	prog := p.ParseProgram()

	checkErrors(t, p)
	if len(prog.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d", len(prog.Statements))
	}

	stmt, ok := prog.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("prog.Statement[0] is not ast.ExpressionStatement. got=%T", prog.Statements[0])
	}

	ident, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("exp not *ast.Identifier. got=%T", stmt.Expression)
	}

	if ident.Value != "foobar" {
		t.Errorf("ident.Value not %s. got=%s", "foobar", ident.Value)
	}

	if ident.TokenLiteral() != "foobar" {
		t.Errorf("ident.TokenLiteral not %s. got=%s", "foobar", ident.TokenLiteral())
	}
}

func checkErrors(t *testing.T, p *Parser) {
	t.Helper()
	errs := p.Errors()
	if len(errs) == 0 {
		return
	}
	t.Errorf("Found %d errors!", len(errs))
	for _, err := range errs {
		t.Errorf("parse error: %s", err)
	}
	t.FailNow()
}
