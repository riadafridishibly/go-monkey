package ast

import (
	"testing"

	"github.com/riadafridishibly/go-monkey/token"
)

func TestString(t *testing.T) {
	prog := &Program{
		Statements: []Statement{
			&LetStatement{
				Token: token.Token{
					Type:    token.LET,
					Literal: "let",
				},
				Name: &Identifier{
					Token: token.Token{
						Type:    token.IDENT,
						Literal: "myVar",
					},
					Value: "myVar",
				},
				Value: &Identifier{
					Token: token.Token{
						Type:    token.IDENT,
						Literal: "anotherVar",
					},
					Value: "anotherVar",
				},
			},
		},
	}

	t.Log(prog.String())

	if prog.String() != "let myVar = anotherVar;" {
		t.Errorf("prog.String() wrong. got=%q", prog.String())
	}
}
