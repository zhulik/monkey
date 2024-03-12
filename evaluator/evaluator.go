package evaluator

import (
	"errors"
	"fmt"

	"github.com/zhulik/monkey/ast"
	obj "github.com/zhulik/monkey/evaluator/object"
)

var (
	ErrParsingError          = errors.New("parsing error")
	ErrUnknownInfixOperator  = errors.New("unknown infix operator")
	ErrUnknownPrefixOperator = errors.New("unknown prefix operator")

	ErrNonBoolCondition = errors.New("condition must be a Boolean")
)

type ReturnValue struct {
	obj.Object
}

type Evaluator struct{}

func New() Evaluator {
	return Evaluator{}
}

func (e Evaluator) Eval(node ast.Node, envs ...obj.EnvGetSetter) (obj.Object, error) { //nolint:cyclop,funlen
	var env obj.EnvGetSetter
	if len(envs) == 0 {
		env = obj.NewEnv()
	} else {
		env = envs[0]
	}

	switch node := node.(type) {
	case *ast.Program:
		return e.evalProgram(node, env)

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
		return e.evalPrefixExpression(node, env)

	case *ast.InfixExpression:
		return e.evalInfixExpression(node, env)

	case *ast.IfExpression:
		return e.evalIfExpression(node, env)

	case *ast.BlockStatement:
		if node != nil && len(node.V) > 0 {
			return e.evalBlockStatement(node, env)
		}

		return obj.NIL, nil

	case *ast.ReturnStatement:
		return e.evalReturnStatement(node, env)

	case *ast.LetStatement:
		return e.evalLetStatement(node, env)

	case *ast.IdentifierExpression:
		return e.evalIdentifierExpression(node, env)

	case *ast.FunctionExpression:
		return obj.Function{Evaluator: e, Function: node, Env: env}, nil

	case *ast.CallExpression:
		return e.evalCallExpression(node, env)

	case *ast.StringExpression:
		return obj.New[obj.String](node.V), nil

	default:
		return nil, fmt.Errorf("%w: unknown node type: %s", ErrParsingError, node.TokenLiteral())
	}
}

func (e Evaluator) evalProgram(node *ast.Program, env obj.EnvGetSetter) (obj.Object, error) {
	var result obj.Object

	var err error

	for _, statement := range node.Statements {
		result, err = e.Eval(statement, env)
		if err != nil {
			return nil, err
		}

		if ret, ok := result.(ReturnValue); ok {
			return ret.Object, nil
		}
	}

	return result, nil
}

func (e Evaluator) evalPrefixExpression(node *ast.PrefixExpression, env obj.EnvGetSetter) (obj.Object, error) {
	value, err := e.Eval(node.V, env)
	if err != nil {
		return nil, err
	}

	switch node.Operator {
	case "-":
		op, cErr := obj.CastOperator[obj.OperatorPrefixMinus](value)
		if cErr != nil {
			return obj.NIL, cErr
		}

		return op.OperatorPrefixMinus() //nolint:wrapcheck
	case "!":
		op, cErr := obj.CastOperator[obj.OperatorBang](value)
		if cErr != nil {
			return obj.NIL, cErr
		}

		return op.OperatorBang() //nolint:wrapcheck

	default:
		return obj.NIL, fmt.Errorf("%w: %s", ErrUnknownPrefixOperator, node.Operator)
	}
}

func (e Evaluator) evalInfixExpression(node *ast.InfixExpression, env obj.EnvGetSetter) (obj.Object, error) { //nolint:funlen,lll,cyclop
	left, err := e.Eval(node.V, env)
	if err != nil {
		return nil, err
	}

	right, err := e.Eval(node.Right, env)
	if err != nil {
		return nil, err
	}

	switch node.Operator {
	case "<":
		op, cErr := obj.CastOperator[obj.OperatorLT](left)
		if cErr != nil {
			return obj.NIL, cErr
		}

		return op.OperatorLT(right) //nolint:wrapcheck

	case "<=":
		op, cErr := obj.CastOperator[obj.OperatorLTE](left)
		if cErr != nil {
			return obj.NIL, cErr
		}

		return op.OperatorLTE(right) //nolint:wrapcheck

	case ">":
		op, cErr := obj.CastOperator[obj.OperatorGT](left)
		if cErr != nil {
			return obj.NIL, cErr
		}

		return op.OperatorGT(right) //nolint:wrapcheck

	case ">=":
		op, cErr := obj.CastOperator[obj.OperatorGTE](left)
		if cErr != nil {
			return obj.NIL, cErr
		}

		return op.OperatorGTE(right) //nolint:wrapcheck

	case "-":
		op, cErr := obj.CastOperator[obj.OperatorMinus](left)
		if cErr != nil {
			return obj.NIL, cErr
		}

		return op.OperatorMinus(right) //nolint:wrapcheck

	case "+":
		op, cErr := obj.CastOperator[obj.OperatorPlus](left)
		if cErr != nil {
			return obj.NIL, cErr
		}

		return op.OperatorPlus(right) //nolint:wrapcheck

	case "*":
		op, cErr := obj.CastOperator[obj.OperatorAsterisk](left)
		if cErr != nil {
			return obj.NIL, cErr
		}

		return op.OperatorAsterisk(right) //nolint:wrapcheck

	case "/":
		op, cErr := obj.CastOperator[obj.OperatorSlash](left)
		if cErr != nil {
			return obj.NIL, cErr
		}

		return op.OperatorSlash(right) //nolint:wrapcheck

	case "==":
		op, cErr := obj.CastOperator[obj.OperatorEQ](left)
		if cErr != nil {
			return obj.NIL, cErr
		}

		return op.OperatorEQ(right) //nolint:wrapcheck

	case "!=":
		op, cErr := obj.CastOperator[obj.OperatorNEQ](left)
		if cErr != nil {
			return obj.NIL, cErr
		}

		return op.OperatorNEQ(right) //nolint:wrapcheck

	default:
		return obj.NIL, fmt.Errorf("%w: %s", ErrUnknownInfixOperator, node.Operator)
	}
}

func (e Evaluator) evalIfExpression(node *ast.IfExpression, env obj.EnvGetSetter) (obj.Object, error) {
	condition, err := e.Eval(node.V, env)
	if err != nil {
		return nil, err
	}

	cond, ok := condition.(obj.Boolean)
	if !ok {
		return obj.NIL, fmt.Errorf("%w, given: %s", ErrNonBoolCondition, condition.TypeName())
	}

	if cond.Value() {
		return e.Eval(node.Then, env)
	}

	return e.Eval(node.Else, env)
}

func (e Evaluator) evalReturnStatement(node *ast.ReturnStatement, env obj.EnvGetSetter) (obj.Object, error) {
	value, err := e.Eval(node.V, env)
	if err != nil {
		return nil, err
	}

	return ReturnValue{value}, nil
}

func (e Evaluator) evalLetStatement(node *ast.LetStatement, env obj.EnvGetSetter) (obj.Object, error) {
	value, err := e.Eval(node.V, env)
	if err != nil {
		return nil, err
	}

	env.Set(node.Name.V, value)

	return value, nil
}

func (e Evaluator) evalIdentifierExpression(node *ast.IdentifierExpression, env obj.EnvGetSetter) (obj.Object, error) {
	return env.Get(node.V) //nolint:wrapcheck
}

func (e Evaluator) evalCallExpression(node *ast.CallExpression, env obj.EnvGetSetter) (obj.Object, error) {
	function, err := e.Eval(node.V, env)
	if err != nil {
		return nil, err
	}

	args := []obj.Object{}

	for _, a := range node.Arguments {
		val, eErr := e.Eval(a, env)
		if eErr != nil {
			return nil, eErr
		}

		args = append(args, val)
	}

	res, err := function.(obj.Function).Call(args...)
	if err != nil {
		return obj.NIL, err //nolint:wrapcheck
	}

	if ret, ok := res.(ReturnValue); ok {
		return ret.Object, nil
	}

	return res, nil
}

func (e Evaluator) evalBlockStatement(node *ast.BlockStatement, env obj.EnvGetSetter) (obj.Object, error) {
	var result obj.Object

	var err error

	for _, statement := range node.V {
		result, err = e.Eval(statement, env)
		if err != nil {
			return nil, err
		}

		if ret, ok := result.(ReturnValue); ok {
			return ret, nil
		}
	}

	return result, nil
}
