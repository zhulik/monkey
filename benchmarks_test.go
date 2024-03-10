package monkey_test

import (
	"math/rand"
	"testing"

	"github.com/samber/lo"
	"github.com/zhulik/monkey/ast"
	"github.com/zhulik/monkey/evaluator"
	obj "github.com/zhulik/monkey/evaluator/object"
	"github.com/zhulik/monkey/lexer"
	"github.com/zhulik/monkey/parser"
)

func random(n int) int {
	return rand.Intn(n) //nolint:gosec
}

func fibNative(n int) int {
	if n < 2 {
		return n
	}

	return fibNative(n-1) + fibNative(n-2)
}

func fibScript(eval evaluator.Evaluator, env obj.EnvGetSetter, program *ast.Program, n int) int {
	env.Set("x", obj.New[obj.Integer](int64(n)))

	val := lo.Must(eval.Eval(program, env))

	return int(val.(obj.Integer).Value()) //nolint:forcetypeassert
}

func BenchmarkFibNative(b *testing.B) {
	for n := 0; n < b.N; n++ {
		fibNative(random(10))
	}
}

func BenchmarkScript(b *testing.B) {
	script := `
let fib = fn(n) {
	if (n < 2) {
		n;
	} else {
		fib(n - 1) + fib(n -  2)
	}
};
fib(x);
`
	lex := lexer.New(script)
	parser := parser.New(lex)

	program := lo.Must(parser.ParseProgram())
	env := obj.NewEnv()

	eval := evaluator.New()

	for n := 0; n < b.N; n++ {
		fibScript(eval, env, program, random(10))
	}
}
