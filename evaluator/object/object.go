package object

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/samber/lo"
)

var (
	NIL   = New[Nil](nil)       //nolint:gochecknoglobals
	TRUE  = New[Boolean](true)  //nolint:gochecknoglobals
	FALSE = New[Boolean](false) //nolint:gochecknoglobals

	ErrUndefinedMethod = errors.New("method is not defined")

	operators = []lo.Tuple2[string, string]{ //nolint:gochecknoglobals
		{A: "<=", B: "LTE"},
		{A: ">=", B: "GTE"},
		{A: "!=", B: "NEQ"},
		{A: "!", B: "Bang"},
		{A: "+", B: "Plus"},
		{A: "-", B: "Minus"},
		{A: "*", B: "Asterisk"},
		{A: "/", B: "Slash"},
		{A: ">", B: "GT"},
		{A: "<", B: "LT"},
		{A: "==", B: "EQ"},
	}
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

func try[T any](fun func() (T, error)) (result T, err error) { //nolint:nonamedreturns
	defer func() {
		if pan := recover(); pan != nil {
			if vErr, ok := pan.(error); ok {
				err = vErr
			} else {
				if strings.Contains(pan.(string), "reflect: Call using") && //nolint:forcetypeassert
					strings.Contains(pan.(string), "as type") { //nolint:forcetypeassert
					err = fmt.Errorf("%w: %w", ErrWronArgumentType, err)

					return
				}

				panic(pan)
			}
		}
	}()

	result, err = fun()

	return
}

func Send(target Object, name string, objects ...Object) (Object, error) {
	method := reflect.ValueOf(target).MethodByName(methodNameToGo(name)) // TODO: Check method defined in the script

	if !method.IsValid() {
		return NIL, fmt.Errorf("%w: %s on %s", ErrUndefinedMethod, name, GetType(target))
	}

	vals, err := try(func() ([]reflect.Value, error) {
		return method.Call(lo.Map(objects, func(item Object, _ int) reflect.Value {
			return reflect.ValueOf(item)
		})), nil
	})
	if err != nil {
		return NIL, err
	}

	if len(vals) != 2 { //nolint:gomnd
		panic("Operators must return exactly 2 values")
	}

	var result Object

	result, ok := vals[0].Interface().(Object)
	if !ok {
		panic("Operator's first return value must be Object, Name:" + name + " Given: " + GetType(vals[0]))
	}

	errOrNil := vals[1].Interface()
	if errOrNil == nil {
		return result, nil
	}

	err, ok = errOrNil.(error)
	if !ok {
		panic("Operator's second return value must be error")
	}

	return result, err
}

func GetType[T any](myvar T) string {
	var res string

	t := reflect.TypeOf(myvar) //nolint:varnamelen
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
		res += "*"
	}

	return res + t.Name()
}

func New[T valuable[V], V any, PT interface {
	*T
	setValue(v V)
}](v V) T {
	res := PT(new(T))
	res.setValue(v)

	return *res
}

type Nil struct {
	BaseObject[any]
}

func (o Nil) TypeName() string {
	return "Nil"
}

func (o Nil) Inspect() string {
	return "nil"
}

func methodNameToGo(name string) string {
	name = capitalize(name)

	for _, v := range operators {
		name = strings.ReplaceAll(name, v.A, v.B)
	}

	return name
}

func capitalize(name string) string {
	r, size := utf8.DecodeRuneInString(name)
	if r == utf8.RuneError {
		return "ERROR"
	}

	return string(unicode.ToUpper(r)) + name[size:]
}
