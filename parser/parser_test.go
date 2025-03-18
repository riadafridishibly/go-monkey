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
