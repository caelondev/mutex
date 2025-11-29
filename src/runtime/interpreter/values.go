package interpreter

type ValueTypes string

const (
	NULL_VALUE   ValueTypes = "null"
	NUMBER_VALUE ValueTypes = "number"
)

type RuntimeValue interface {
	Type() ValueTypes
}

type NullValue struct{}

func (n *NullValue) Type() ValueTypes {
	return NULL_VALUE
}

type NumberValue struct {
	Value float64
}

func (n *NumberValue) Type() ValueTypes {
	return NUMBER_VALUE
}
