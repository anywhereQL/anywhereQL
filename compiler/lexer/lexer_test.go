package lexer

import (
	"testing"

	"github.com/anywhereQL/anywhereQL/common/token"
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
