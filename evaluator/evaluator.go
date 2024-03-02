package evaluator

import (
	"errors"
	"fmt"

	"github.com/zhulik/monkey/ast"
	obj "github.com/zhulik/monkey/evaluator/object"
)

var ErrParsingError = errors.New("parsing error")

type ReturnValue struct {
	obj.Object
}

type Evaluator struct{}

func (e Evaluator) Eval(node ast.Node, envs ...obj.EnvGetSetter) (obj.Object, error) { //nolint:cyclop,funlen
	var env obj.EnvGetSetter
	if len(envs) == 0 {
		env = obj.NewEnv()
	} else {
		env = envs[0]
	}

	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(e, node, env)

	case *ast.ExpressionStatement:
		return e.Eval(node.V, env)

	case *ast.BooleanExpression:
		if node.V {
			return obj.TRUE, nil
		} else {
			return obj.FALSE, nil
		}

	case *ast.NilExpression:
		return obj.NIL, nil

	case *ast.IntegerExpression:
		return obj.New[obj.Integer](node.V), nil

	case *ast.PrefixExpression:
		return evalPrefixExpression(e, node, env)

	case *ast.InfixExpression:
		return evalInfixExpression(e, node, env)

	case *ast.IfExpression:
		return evalIfExpression(e, node, env)

	case *ast.BlockStatement:
		if node != nil && len(node.V) > 0 {
			return evalBlockStatement(e, node, env)
		}

		return obj.NIL, nil

	case *ast.ReturnStatement:
		return evalReturnStatement(e, node, env)

	case *ast.LetStatement:
		return evalLetStatement(e, node, env)

	case *ast.IdentifierExpression:
		return evalIdentifierExpression(node, env)

	case *ast.FunctionExpression:
		return obj.Function{Evaluator: e, Function: node, Env: env}, nil

	case *ast.CallExpression:
		return evalCallExpression(e, node, env)

	default:
		return nil, fmt.Errorf("%w: unknown node type: %s", ErrParsingError, obj.GetType(node))
	}
}

func New() Evaluator {
	return Evaluator{}
}

// func evalStatements(e Evaluator, statements []ast.Statement, env obj.EnvGetSetter) (obj.Object, error) {
// 	var result obj.Object

// 	var err error
// 	for _, stmt := range statements {
// 		result, err = e.Eval(stmt, env)
// 		if err != nil {
// 			return nil, err
// 		}

// 		if ret, ok := result.(ReturnValue); ok {
// 			return ret.Object, nil
// 		}
// 	}

// 	return result, nil
// }

func evalProgram(eval Evaluator, node *ast.Program, env obj.EnvGetSetter) (obj.Object, error) {
	var result obj.Object

	var err error

	for _, statement := range node.Statements {
		result, err = eval.Eval(statement, env)
		if err != nil {
			return nil, err
		}

		if ret, ok := result.(ReturnValue); ok {
			return ret.Object, nil
		}
	}

	return result, nil
}

func evalPrefixExpression(eval Evaluator, node *ast.PrefixExpression, env obj.EnvGetSetter) (obj.Object, error) {
	value, err := eval.Eval(node.V, env)
	if err != nil {
		return nil, err
	}

	operator := node.Operator
	if operator == "-" {
		operator = "PrefixMinus"
	}

	result, err := obj.Send(value, "operator"+operator)
	if err != nil {
		return nil, fmt.Errorf("send error: %w", err)
	}

	return result, nil
}

func evalInfixExpression(eval Evaluator, node *ast.InfixExpression, env obj.EnvGetSetter) (obj.Object, error) {
	left, err := eval.Eval(node.V, env)
	if err != nil {
		return nil, err
	}

	right, err := eval.Eval(node.Right, env)
	if err != nil {
		return nil, err
	}

	result, err := obj.Send(left, "operator"+node.Operator, right)
	if err != nil {
		return nil, fmt.Errorf("send error: %w", err)
	}

	return result, nil
}

func evalIfExpression(eval Evaluator, node *ast.IfExpression, env obj.EnvGetSetter) (obj.Object, error) {
	condition, err := eval.Eval(node.V, env)
	if err != nil {
		return nil, err
	}

	cond, ok := condition.(obj.Boolean)
	if !ok {
		panic("condition must be bool")
	}

	if cond.Value() {
		return eval.Eval(node.Then, env)
	}

	return eval.Eval(node.Else, env)
}

func evalReturnStatement(eval Evaluator, node *ast.ReturnStatement, env obj.EnvGetSetter) (obj.Object, error) {
	value, err := eval.Eval(node.V, env)
	if err != nil {
		return nil, err
	}

	return ReturnValue{value}, nil
}

func evalLetStatement(eval Evaluator, node *ast.LetStatement, env obj.EnvGetSetter) (obj.Object, error) {
	value, err := eval.Eval(node.V, env)
	if err != nil {
		return nil, err
	}

	env.Set(node.Name.V, value)

	return value, nil
}

func evalIdentifierExpression(node *ast.IdentifierExpression, env obj.EnvGetSetter) (obj.Object, error) {
	return env.Get(node.V) //nolint:wrapcheck
}

func evalCallExpression(eval Evaluator, node *ast.CallExpression, env obj.EnvGetSetter) (obj.Object, error) {
	function, err := eval.Eval(node.V, env)
	if err != nil {
		return nil, err
	}

	args := []obj.Object{}

	for _, a := range node.Arguments {
		val, eErr := eval.Eval(a, env)
		if eErr != nil {
			return nil, eErr
		}

		args = append(args, val)
	}

	return function.(obj.Function).Call(args...) //nolint:forcetypeassert,wrapcheck
}

func evalBlockStatement(eval Evaluator, node *ast.BlockStatement, env obj.EnvGetSetter) (obj.Object, error) {
	var result obj.Object

	var err error

	for _, statement := range node.V {
		result, err = eval.Eval(statement, env)
		if err != nil {
			return nil, err
		}

		if ret, ok := result.(ReturnValue); ok {
			return ret, nil
		}
	}

	return result, nil
}
