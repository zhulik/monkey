package object

type Nil struct {
	BaseObject[any]
}

func (o Nil) TypeName() string {
	return "Nil"
}

func (o Nil) Inspect() string {
	return "nil"
}

func (o Nil) OperatorEQ(other Object) (Object, error) {
	_, ok := other.(Nil)
	if !ok {
		return NIL, ErrWronArgumentType
	}

	return TRUE, nil
}

func (o Nil) OperatorNEQ(other Object) (Object, error) {
	_, ok := other.(Nil)
	if !ok {
		return NIL, ErrWronArgumentType
	}

	return FALSE, nil
}
