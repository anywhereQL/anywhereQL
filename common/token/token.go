package token

import (
	"strings"

	"github.com/anywhereQL/anywhereQL/common/value"
)

type Type int

const (
	UNKNOWN Type = -1
	ERROR   Type = iota
	EOS
	IDENT
	NUMBER
	STRING

	K_SELECT
	K_NULL
	K_FROM
	K_TRUE
	K_FALSE

	S_PLUS
	S_MINUS
	S_ASTERISK
	S_SOLIDAS
	S_PERCENT
	S_SEMICOLON
	S_LPAREN
	S_RPAREN
	S_COMMA
	S_QUOTE
	S_DQUOTE
	S_PERIOD
)

func (t Type) String() string {
	switch t {
	case UNKNOWN:
		return "Unknwon Token"
	case ERROR:
		return "Error Token"
	case EOS:
		return "End of SQL"
	case IDENT:
		return "IDENT Token"
	case NUMBER:
		return "NUMBER Token"
	case STRING:
		return "STRING Token"

	case K_SELECT:
		return "Keyword (SELECT)"
	case K_NULL:
		return "Keyword (NULL)"
	case K_FROM:
		return "Keyword (FROM)"
	case K_TRUE:
		return "Keyword (TRUE)"
	case K_FALSE:
		return "Keyword (FALSE)"

	case S_PLUS:
		return "Symbol (+)"
	case S_MINUS:
		return "Symbol (-)"
	case S_ASTERISK:
		return "Symbol (*)"
	case S_SOLIDAS:
		return "Symbol (/)"
	case S_PERCENT:
		return "Symbol (%)"
	case S_SEMICOLON:
		return "Symbol (;)"
	case S_LPAREN:
		return "Symbol (()"
	case S_RPAREN:
		return "Symbol ())"
	case S_COMMA:
		return "Symbol (,)"
	case S_QUOTE:
		return "Symbol (')"
	case S_DQUOTE:
		return "Symbol (\")"
	case S_PERIOD:
		return "Symbol (.)"

	default:
		return "Error!! Unknown Token Type"
	}
}

type Token struct {
	Type    Type
	Literal string
	Value   value.Value
}

type Tokens []Token

func (t Tokens) GetN(n int) Token {
	if len(t) <= n {
		return Token{Type: ERROR}
	}
	return t[n]
}

func LookupKeyword(s string) (bool, Type) {
	switch strings.ToUpper(s) {
	case "SELECT":
		return true, K_SELECT
	case "NULL":
		return true, K_NULL
	case "FROM":
		return true, K_FROM
	case "TRUE":
		return true, K_TRUE
	case "FALSE":
		return true, K_FALSE
	default:
		return false, UNKNOWN
	}
}
