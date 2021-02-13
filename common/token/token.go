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
	K_AND
	K_OR
	K_IS
	K_NOT
	K_ISNULL
	K_NOTNULL
	K_CAST
	K_INT
	K_INTEGER
	K_FLOAT
	K_DOUBLE
	K_STRING
	K_AS

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
	S_EQUAL
	S_NOT_EQUAL
	S_LESS_THAN
	S_LESS_THAN_EQUAL
	S_GREATER_THAN
	S_GREATER_THAN_EQUAL
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
	case K_AND:
		return "Keyword (AND)"
	case K_OR:
		return "Keyword (OR)"
	case K_IS:
		return "Keyword (IS)"
	case K_NOT:
		return "Keyword (NOT)"
	case K_ISNULL:
		return "Keyword (ISNULL)"
	case K_NOTNULL:
		return "Keyword (NOTNULL)"
	case K_CAST:
		return "Keyword (CAST)"
	case K_INT:
		return "Keyword (INT)"
	case K_INTEGER:
		return "Keyword (INTEGER)"
	case K_FLOAT:
		return "Keyword (FLOAT)"
	case K_DOUBLE:
		return "Keyword (DOUBLE)"
	case K_STRING:
		return "Keyword (STRING)"
	case K_AS:
		return "Keyword (AS)"

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
	case S_EQUAL:
		return "Symbol (=)"
	case S_NOT_EQUAL:
		return "Symbol (<>)"
	case S_LESS_THAN:
		return "Symbol (<)"
	case S_LESS_THAN_EQUAL:
		return "Symbol (<=)"
	case S_GREATER_THAN:
		return "Symbol (>)"
	case S_GREATER_THAN_EQUAL:
		return "Symbol (>=)"

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
	case "AND":
		return true, K_AND
	case "OR":
		return true, K_OR
	case "IS":
		return true, K_IS
	case "NOT":
		return true, K_NOT
	case "ISNULL":
		return true, K_ISNULL
	case "NOTNULL":
		return true, K_NOTNULL
	case "CAST":
		return true, K_CAST
	case "INT", "INTEGER":
		return true, K_INT
	case "FLOAT", "DOUBLE":
		return true, K_DOUBLE
	case "STRING":
		return true, K_STRING
	case "AS":
		return true, K_AS

	default:
		return false, UNKNOWN
	}
}
