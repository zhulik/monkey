package object

import "fmt"

type String struct {
	BaseObject[string]
}

func (s String) TypeName() string {
	return "String"
}

func (s String) Inspect() string {
	return fmt.Sprintf(`"%s"`, s.value)
}

func (s String) OperatorPlus(other Object) (Object, error) {
	otherInt, ok := other.(String)
	if !ok {
		return NIL, ErrWronArgumentType
	}

	s.value += otherInt.value

	return s, nil
}

func (s String) OperatorEQ(other Object) (Object, error) {
	otherInt, ok := other.(String)
	if !ok {
		return NIL, ErrWronArgumentType
	}

	return ToBoolean(s.value == otherInt.value), nil
}

func (s String) OperatorNEQ(other Object) (Object, error) {
	otherInt, ok := other.(String)
	if !ok {
		return NIL, ErrWronArgumentType
	}

	return ToBoolean(s.value != otherInt.value), nil
}
