package result

type ResultType int

const (
	_ ResultType = iota
	Integral
)

type Value struct {
	Type     ResultType
	Integral int64
}
