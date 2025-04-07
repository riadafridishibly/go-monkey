package parser

import (
	"fmt"
	"testing"

	"github.com/k0kubun/pp/v3"
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
	checkParserErrors(t, p)

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

	checkParserErrors(t, p)
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

func TestParsingPrefixExpression(t *testing.T) {
	prefixTests := []struct {
		input    string
		operator string
		intValue int64
	}{
		{"!5;", "!", 5},
		{"-15", "-", 15},
	}

	for _, tt := range prefixTests {
		l := lexer.New(tt.input)
		p := New(l)
		prog := p.ParseProgram()
		checkParserErrors(t, p)

		if len(prog.Statements) != 1 {
			t.Fatalf("prog.Statements does not contain %d statements. got=%d\n",
				1, len(prog.Statements))
		}

		stmt, ok := prog.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("prog.Statement[0] is not ast.ExpressionStatement. got=%T",
				prog.Statements[0],
			)
		}

		expr, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("stmt is not ast.PrefixExpression. got=%T", stmt.Expression)
		}

		if expr.Operator != tt.operator {
			t.Fatalf("expr.Operator is not %s. got %s", tt.operator, expr.Operator)
		}

		if !testIntegerLiteral(t, expr.Right, tt.intValue) {
			return
		}

	}
}

func TestParseInfixExpression(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  int64
		operator   string
		rightValue int64
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
	}

	for _, tt := range infixTests {
		l := lexer.New(tt.input)
		p := New(l)
		prog := p.ParseProgram()
		checkParserErrors(t, p)

		if len(prog.Statements) != 1 {
			t.Fatalf("expected %d statements. got=%d", 1, len(prog.Statements))
		}

		// TODO: What does an ExpressionStatement look like anyway?
		stmt, ok := prog.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("expected expresson statement. got=%T", prog.Statements[0])
		}
		expr, ok := stmt.Expression.(*ast.InfixExpression)
		if !ok {
			t.Fatalf("expected ast.InfixExpression. got=%T", stmt.Expression)
		}

		if expr.Operator != tt.operator {
			t.Fatalf("expected operator %q but got %q", tt.operator, expr.Operator)
		}
		if !testIntegerLiteral(t, expr.Right, tt.rightValue) {
			return
		}
	}
}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	t.Helper()
	integer, ok := il.(*ast.InetegerLiteral)
	if !ok {
		t.Errorf("expected *ast.IntegerLiteral. got=%T", il)
		return false
	}
	if integer.Value != value {
		t.Errorf("expected Value %d. got=%d", value, integer.Value)
		return false
	}
	if integer.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("expected TokenLiteral %d, got %s", value, integer.TokenLiteral())
		return false
	}
	return true
}

func checkParserErrors(t *testing.T, p *Parser) {
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

func TestParseABC(t *testing.T) {
	l := lexer.New("a + b * c")
	p := New(l)
	prog := p.ParseProgram()

	pp.Default.SetColoringEnabled(false)
	pp.Println(prog)
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	var testCases = []struct {
		input    string
		expected string
	}{
		{
			input:    "-a * b",
			expected: "((-a) * b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",
		},
		{
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		},
		{
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
	}

	for _, tc := range testCases {
		l := lexer.New(tc.input)
		p := New(l)
		prog := p.ParseProgram()
		checkParserErrors(t, p)

		// This string method will wrap things in parenthesis
		actual := prog.String()

		if actual != tc.expected {
			t.Errorf("failed! expected = %q but got = %q", tc.expected, actual)
		}
	}
}
