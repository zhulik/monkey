package object

import (
	"strconv"
)

type Boolean struct {
	BaseObject[bool]
}

func (o Boolean) TypeName() string {
	return "Boolean"
}

func (o Boolean) Inspect() string {
	return strconv.FormatBool(o.value)
}

func (o Boolean) OperatorBang() (Object, error) {
	if o.value {
		return FALSE, nil
	}

	return TRUE, nil
}

func (o Boolean) OperatorEQ(other Object) (Object, error) {
	otherBool, ok := other.(Boolean)
	if !ok {
		return NIL, ErrWronArgumentType
	}

	return New[Boolean](o.value == otherBool.value), nil
}

func (o Boolean) OperatorNEQ(other Object) (Object, error) {
	otherBool, ok := other.(Boolean)
	if !ok {
		return NIL, ErrWronArgumentType
	}

	return New[Boolean](o.value != otherBool.value), nil
}
