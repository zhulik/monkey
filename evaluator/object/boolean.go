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

func (o Boolean) OperatorBang() (Boolean, error) {
	if o.value {
		return FALSE, nil
	}

	return TRUE, nil
}

func (o Boolean) OperatorEQ(other Boolean) (Boolean, error) {
	return New[Boolean](o.value == other.value), nil
}

func (o Boolean) OperatorNEQ(other Boolean) (Boolean, error) {
	return New[Boolean](o.value != other.value), nil
}
