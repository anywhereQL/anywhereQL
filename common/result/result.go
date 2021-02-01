package result

type ResultType int

const (
	_ ResultType = iota
	Integral
	Float
	Decimal
)

type Value struct {
	Type     ResultType
	Integral int64
	Float    float64
	PartF    int64
	PartI    int64
	FDigit   int
}
