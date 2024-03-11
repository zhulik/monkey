package object

var (
	NIL   = New[Nil](nil)       //nolint:gochecknoglobals
	TRUE  = New[Boolean](true)  //nolint:gochecknoglobals
	FALSE = New[Boolean](false) //nolint:gochecknoglobals
)

type Object interface {
	TypeName() string
	Inspect() string
}

type valuable[V any] interface {
	Value() V
}

type BaseObject[T any] struct {
	value T
}

func (i BaseObject[T]) Value() T {
	return i.value
}

func (i *BaseObject[T]) setValue(v T) { //nolint:unused
	i.value = v
}

func New[T valuable[V], V any, PT interface {
	*T
	setValue(v V)
}](v V) T {
	res := PT(new(T))
	res.setValue(v)

	return *res
}
