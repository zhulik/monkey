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

func (o Integer) OperatorPlus(other Object) (Object, error) {
	otherInt, ok := other.(Integer)
	if !ok {
		return NIL, ErrWronArgumentType
	}

	return New[Integer](o.value + otherInt.value), nil
}

func (o Integer) OperatorPrefixMinus() (Object, error) {
	return New[Integer](-o.value), nil
}

func (o Integer) OperatorMinus(other Object) (Object, error) {
	otherInt, ok := other.(Integer)
	if !ok {
		return NIL, ErrWronArgumentType
	}

	return New[Integer](o.value - otherInt.value), nil
}

func (o Integer) OperatorAsterisk(other Object) (Object, error) {
	otherInt, ok := other.(Integer)
	if !ok {
		return NIL, ErrWronArgumentType
	}

	return New[Integer](o.value * otherInt.value), nil
}

func (o Integer) OperatorSlash(other Object) (Object, error) {
	otherInt, ok := other.(Integer)
	if !ok {
		return NIL, ErrWronArgumentType
	}

	if otherInt.value == 0 {
		return Integer{}, ErrDevisionByZero
	}

	return New[Integer](o.value / otherInt.value), nil
}

func (o Integer) OperatorGT(other Object) (Object, error) {
	otherInt, ok := other.(Integer)
	if !ok {
		return NIL, ErrWronArgumentType
	}

	return New[Boolean](o.value > otherInt.value), nil
}

func (o Integer) OperatorGTE(other Object) (Object, error) {
	otherInt, ok := other.(Integer)
	if !ok {
		return NIL, ErrWronArgumentType
	}

	return New[Boolean](o.value >= otherInt.value), nil
}

func (o Integer) OperatorLT(other Object) (Object, error) {
	otherInt, ok := other.(Integer)
	if !ok {
		return NIL, ErrWronArgumentType
	}

	return New[Boolean](o.value < otherInt.value), nil
}

func (o Integer) OperatorLTE(other Object) (Object, error) {
	otherInt, ok := other.(Integer)
	if !ok {
		return NIL, ErrWronArgumentType
	}

	return New[Boolean](o.value <= otherInt.value), nil
}

func (o Integer) OperatorEQ(other Object) (Object, error) {
	otherInt, ok := other.(Integer)
	if !ok {
		return NIL, ErrWronArgumentType
	}

	return New[Boolean](o.value == otherInt.value), nil
}

func (o Integer) OperatorNEQ(other Object) (Object, error) {
	otherInt, ok := other.(Integer)
	if !ok {
		return NIL, ErrWronArgumentType
	}

	return New[Boolean](o.value != otherInt.value), nil
}
