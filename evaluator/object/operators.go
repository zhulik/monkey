package object

import (
	"errors"

	"github.com/samber/lo"
)

var ErrUndefinedMethod = errors.New("method is not defined")

// Comparison operators.
type OperatorLT interface {
	OperatorLT(other Object) (Object, error)
}

type OperatorLTE interface {
	OperatorLTE(other Object) (Object, error)
}

type OperatorGT interface {
	OperatorGT(other Object) (Object, error)
}

type OperatorGTE interface {
	OperatorGTE(other Object) (Object, error)
}

type OperatorEQ interface {
	OperatorEQ(other Object) (Object, error)
}

type OperatorNEQ interface {
	OperatorNEQ(other Object) (Object, error)
}

// Prefix operators.
type OperatorPrefixMinus interface {
	OperatorPrefixMinus() (Object, error)
}

type OperatorBang interface {
	OperatorBang() (Object, error)
}

// Infix operators.
type OperatorMinus interface {
	OperatorMinus(other Object) (Object, error)
}

type OperatorPlus interface {
	OperatorPlus(other Object) (Object, error)
}

type OperatorAsterisk interface {
	OperatorAsterisk(other Object) (Object, error)
}

type OperatorSlash interface {
	OperatorSlash(other Object) (Object, error)
}

func CastOperator[O any](obj Object) (O, error) {
	op, ok := obj.(O)
	if !ok {
		return lo.Empty[O](), ErrUndefinedMethod
	}

	return op, nil
}
