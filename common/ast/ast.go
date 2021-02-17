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
	Column          *Column
	Cast            *Cast
	Case            *Case
	Between         *Between
}

type Between struct {
	Not   bool
	Src   *Expression
	Begin *Expression
	End   *Expression
}

type Case struct {
	Value      *Expression
	CaseValues []CaseCondition
	ElseValue  *Expression
}

type CaseCondition struct {
	Condition *Expression
	Result    *Expression
}

type Type int

const (
	T_INT Type = iota
	T_FLOAT
	T_STRING
)

type Cast struct {
	Expr *Expression
	Type Type
}

type Literal struct {
	Numeric *Numeric
	String  *String
	Bool    *Bool
	NULL    bool
}

type Bool struct {
	True    bool
	False   bool
	Unknown bool
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
	B_EQUAL
	B_NOT_EQUAL
	B_GREATER_THAN
	B_GREATER_THAN_EQUAL
	B_LESS_THAN
	B_LESS_THAN_EQUAL
	B_AND
	B_OR

	U_PLUS
	U_MINUS
	U_NOT
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
	case B_EQUAL:
		return "="
	case B_NOT_EQUAL:
		return "<>"
	case B_GREATER_THAN:
		return ">"
	case B_GREATER_THAN_EQUAL:
		return ">="
	case B_LESS_THAN:
		return "<"
	case B_LESS_THAN_EQUAL:
		return "<="
	case B_AND:
		return "AND"
	case B_OR:
		return "OR"

	case U_PLUS:
		return "+"
	case U_MINUS:
		return "-"
	case U_NOT:
		return "NOT"
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
	Table  Table
}

type Table struct {
	ID     string
	Table  string
	DB     string
	Schema string
}
