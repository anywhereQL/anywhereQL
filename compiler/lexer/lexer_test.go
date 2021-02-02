package lexer

import (
	"testing"

	"github.com/anywhereQL/anywhereQL/common/token"
	"github.com/anywhereQL/anywhereQL/common/value"
)

func TestLexer(t *testing.T) {
	testCases := []struct {
		input    string
		expected token.Tokens
		isError  bool
	}{
		{
			input: "SELECT",
			expected: token.Tokens{
				{
					Type:    token.K_SELECT,
					Literal: "SELECT",
				},
				{
					Type: token.EOS,
				},
			},
			isError: false,
		},
		{
			input: "SELECT 1",
			expected: token.Tokens{
				{
					Type:    token.K_SELECT,
					Literal: "SELECT",
				},
				{
					Type:    token.NUMBER,
					Literal: "1",
					Value: value.Value{
						Type: value.INTEGER,
						Int:  1,
					},
				},
				{
					Type: token.EOS,
				},
			},
			isError: false,
		},
		{
			input: "SELECT 1 + 2",
			expected: token.Tokens{
				{
					Type:    token.K_SELECT,
					Literal: "SELECT",
				},
				{
					Type:    token.NUMBER,
					Literal: "1",
					Value: value.Value{
						Type: value.INTEGER,
						Int:  1,
					},
				},
				{
					Type:    token.S_PLUS,
					Literal: "+",
				},
				{
					Type:    token.NUMBER,
					Literal: "2",
					Value: value.Value{
						Type: value.INTEGER,
						Int:  2,
					},
				},
				{
					Type: token.EOS,
				},
			},
			isError: false,
		},
		{
			input: "SELECT ABS(1 + 2);",
			expected: token.Tokens{
				{
					Type:    token.K_SELECT,
					Literal: "SELECT",
				},
				{
					Type:    token.IDENT,
					Literal: "ABS",
				},
				{
					Type:    token.S_LPAREN,
					Literal: "(",
				},
				{
					Type:    token.NUMBER,
					Literal: "1",
					Value: value.Value{
						Type: value.INTEGER,
						Int:  1,
					},
				},
				{
					Type:    token.S_PLUS,
					Literal: "+",
				},
				{
					Type:    token.NUMBER,
					Literal: "2",
					Value: value.Value{
						Type: value.INTEGER,
						Int:  2,
					},
				},
				{
					Type:    token.S_RPAREN,
					Literal: ")",
				},
				{
					Type:    token.S_SEMICOLON,
					Literal: ";",
				},
				{
					Type: token.EOS,
				},
			},
			isError: false,
		},
	}

	for tn, tc := range testCases {
		tokens := Lex(tc.input)
		if len(tc.expected) != len(tokens) {
			t.Fatalf("[%d] Length is mismatcg", tn)
		}
		for n, tk := range tokens {
			expected := tc.expected.GetN(n)
			if expected.Type != tk.Type {
				t.Fatalf("[%d] Token Type is mismatch %s:%s", tn, expected.Type, tk.Type)
			}
			if expected.Literal != tk.Literal {
				t.Fatalf("[%d] Token Literal is mismatch", tn)
			}
		}
	}
}
