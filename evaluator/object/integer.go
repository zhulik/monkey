package object

import (
	"errors"
	"strconv"
)

var (
	ErrDevisionByZero   = errors.New("division by zero")
	ErrWronArgumentType = errors.New("argument type is wrong")
)

type Integer struct {
	BaseObject[int64]
}

func (o Integer) TypeName() string {
	return "Integer"
}

func (o Integer) Inspect() string {
	return strconv.FormatInt(o.value, 10)
}

func (o Integer) OperatorPlus(other Integer) (Integer, error) {
	return New[Integer](o.value + other.value), nil
}

func (o Integer) OperatorPrefixMinus() (Integer, error) {
	return New[Integer](-o.value), nil
}

func (o Integer) OperatorMinus(other Integer) (Integer, error) {
	return New[Integer](o.value - other.value), nil
}

func (o Integer) OperatorAsterisk(other Integer) (Integer, error) {
	return New[Integer](o.value * other.value), nil
}

func (o Integer) OperatorSlash(other Integer) (Integer, error) {
	if other.value == 0 {
		return Integer{}, ErrDevisionByZero
	}

	return New[Integer](o.value / other.value), nil
}

func (o Integer) OperatorGT(other Integer) (Boolean, error) {
	return New[Boolean](o.value > other.value), nil
}

func (o Integer) OperatorGTE(other Integer) (Boolean, error) {
	return New[Boolean](o.value >= other.value), nil
}

func (o Integer) OperatorLT(other Integer) (Boolean, error) {
	return New[Boolean](o.value < other.value), nil
}

func (o Integer) OperatorLTE(other Integer) (Boolean, error) {
	return New[Boolean](o.value <= other.value), nil
}

func (o Integer) OperatorEQ(other Integer) (Boolean, error) {
	return New[Boolean](o.value == other.value), nil
}

func (o Integer) OperatorNEQ(other Integer) (Boolean, error) {
	return New[Boolean](o.value != other.value), nil
}
