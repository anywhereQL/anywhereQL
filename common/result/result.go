package result

type ResultType int

const (
	_ ResultType = iota
	Integral
	Float
	Decimal
	String
)

type Value struct {
	Type     ResultType
	Integral int64
	Float    float64
	String   string

	PartF  int64
	PartI  int64
	FDigit int
}
