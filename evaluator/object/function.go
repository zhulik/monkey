package object

import (
	"github.com/zhulik/monkey/ast"
)

type Evaluator interface {
	Eval(node ast.Node, envs ...EnvGetSetter) (Object, error)
}

type Function struct {
	Evaluator Evaluator
	Function  *ast.FunctionExpression
	Env       EnvGetSetter
}

func (o Function) TypeName() string {
	return "Function"
}

func (o Function) Inspect() string {
	return o.Function.String()
}

func (o Function) Call(args ...Object) (Object, error) {
	return o.Evaluator.Eval(o.Function.V, o.addArgsToEnv(args)) //nolint:wrapcheck
}

func (o Function) addArgsToEnv(args []Object) EnvGetSetter {
	closure := NewEnv(o.Env)
	for i, val := range args {
		closure.Set(o.Function.Arguments[i].V, val)
	}

	return closure
}
