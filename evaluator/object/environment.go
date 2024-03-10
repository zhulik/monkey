package object

import (
	"errors"
	"fmt"
)

var ErrUnknownIdentifier = errors.New("identifier is unknown")

type EnvGetter interface {
	Get(name string) (Object, error)
}

type EnvSetter interface {
	Set(name string, val Object) Object
}

type EnvGetSetter interface {
	EnvGetter
	EnvSetter
}

type Env struct {
	store  map[string]Object
	parent EnvGetter
}

func NewEnv(parents ...EnvGetter) *Env {
	var parent EnvGetter
	if len(parents) > 0 {
		parent = parents[0]
	}

	return &Env{
		store:  map[string]Object{},
		parent: parent,
	}
}

func (e Env) Get(name string) (Object, error) {
	val, ok := e.store[name]

	if !ok {
		if e.parent == nil {
			return nil, fmt.Errorf("%w: %s", ErrUnknownIdentifier, name)
		}

		return e.parent.Get(name) //nolint:wrapcheck
	}

	return val, nil
}

func (e *Env) Set(name string, object Object) Object {
	e.store[name] = object

	return object
}
