package ast

type AST struct {
	SQL []SQL
}

type SQL struct {
	SELECTStatement *SELECTStatement
}

type SELECTStatement struct {
	SELECT *SELECTClause
	FROM   *FROMClause
}

type SELECTClause struct {
	SelectColumns []SelectColumn
}

type FROMClause struct {
	Table *Table
}

type SelectColumn struct {
	Expression *Expression
}

type Expression struct {
	Literal         *Literal
	UnaryOperation  *UnaryOpe
	BinaryOperation *BinaryOpe
	FunctionCall    *FunctionCall
}

type Literal struct {
	Numeric *Numeric
	String  *String
}

type String struct {
	Value string
}

type NumericType int

const (
	_ NumericType = iota
	N_INT
	N_FLOAT
)

type Numeric struct {
	Type     NumericType
	Integral int64
	Float    float64
	PartI    int64
	PartF    int64
	FDigit   int
}

type OperatorType int

const (
	_ OperatorType = iota
	B_PLUS
	B_MINUS
	B_ASTERISK
	B_SOLIDAS
	B_PERCENT

	U_PLUS
	U_MINUS
)

func (o OperatorType) String() string {
	switch o {
	case B_PLUS:
		return "+"
	case B_MINUS:
		return "-"
	case B_ASTERISK:
		return "*"
	case B_SOLIDAS:
		return "/"
	case B_PERCENT:
		return "%"

	case U_PLUS:
		return "+"
	case U_MINUS:
		return "-"
	default:
		return "Unknown Operator"
	}
}

type BinaryOpe struct {
	Operator OperatorType
	Left     *Expression
	Right    *Expression
}

type UnaryOpe struct {
	Operator OperatorType
	Expr     *Expression
}

type FunctionCall struct {
	Name string
	Args []Expression
}

type Column struct {
	Column string
	Table  string
	DB     string
	Schema string
}

type Table struct {
	Table  string
	DB     string
	Schema string
	Aslias string
}
