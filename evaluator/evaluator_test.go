package evaluator_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/samber/lo"
	"github.com/zhulik/monkey/evaluator"
	obj "github.com/zhulik/monkey/evaluator/object"
	"github.com/zhulik/monkey/lexer"
	"github.com/zhulik/monkey/parser"
)

func eval(str string) (obj.Object, error) {
	lex := lexer.New(str)
	par := parser.New(lex)
	evaluator := evaluator.New()
	ast := lo.Must(par.ParseProgram())

	return evaluator.Eval(ast)
}

var _ = Describe("Evaluator", func() {
	Describe(".Eval", func() {
		Context("when program is correct", func() {
			cases := map[string]string{
				"10":    "10",
				"123":   "123",
				"true":  "true",
				"false": "false",
				"nil":   "nil",

				"!!true":  "true",
				"!true":   "false",
				"!!false": "false",
				"!false":  "true",

				"-1": "-1",

				"1 + 1": "2",
				"1 - 1": "0",
				"2 * 2": "4",
				"2 / 2": "1",

				"1 == 1": "true",
				"1 != 1": "false",

				"1 > 1": "false",
				"1 < 1": "false",

				"1 >= 1": "true",
				"1 <= 1": "true",

				"true == true":   "true",
				"false == false": "true",
				"true == false":  "false",
				"true != false":  "true",

				"if (true) { 1 } else { 0 }": "1",

				"return 10;":       "10",
				"1; return 10; 1;": "10",
				"if (10 > 1) { if (1 < 10) { return 10; 2 }; 1; }": "10",
				"if (10 < 1) {  } else { 1 }":                      "1",
				"if (10 > 1) { 10 } else { true + true }":          "10",

				"let a = 10; let b = 10; b;":                                                            "10",
				"let add = fn(a, b){ a + b; }":                                                          "fn(a, b) { (a + b) }",
				"let add = fn(a, b){ a + b; }; add(2, 2)":                                               "4",
				"let i = 1; let add = fn(a){ a + i; }; add(2)":                                          "3",
				"let apply = fn(a, b){ b(a) }; apply(2, fn(a) { a + 1 })":                               "3",
				"let check = fn(a){ a == 10 }; check(10)":                                               "true",
				"let check = fn(a){ fn(b) { a == b } }; check(10)(10)":                                  "true",
				"let a = fn() { 1 }; let b = fn() { a(); }; b()":                                        "1",
				"let fib = fn(n) { if (n < 2) { return n; } return fib(n - 1) + fib(n -  2); }; fib(2)": "1",

				`"foo bar"`:     `"foo bar"`,
				`"foo" + "bar"`: `"foobar"`,
			}

			for input, output := range cases {
				Context("when input="+input, func() {
					It("returns "+output, func() {
						result, err := eval(input)
						Expect(err).ToNot(HaveOccurred())
						Expect(result.Inspect()).To(Equal(output))
					})
				})
			}
		})

		Context("when program is incorrect", func() {
			cases := map[string]error{
				"1/0": obj.ErrDevisionByZero,

				"1 > true":  obj.ErrWronArgumentType,
				"1 >= true": obj.ErrWronArgumentType,
				"1 < true":  obj.ErrWronArgumentType,
				"1 <= true": obj.ErrWronArgumentType,
				"1 == true": obj.ErrWronArgumentType,
				"1 != true": obj.ErrWronArgumentType,

				"true > 1":  obj.ErrUndefinedMethod,
				"true >= 1": obj.ErrUndefinedMethod,
				"true < 1":  obj.ErrUndefinedMethod,
				"true <= 1": obj.ErrUndefinedMethod,
				"true == 1": obj.ErrWronArgumentType,
				"true != 1": obj.ErrWronArgumentType,

				"if (10 < 1) { 10 } else { true + true }": obj.ErrUndefinedMethod,
			}

			for input, resultErr := range cases {
				Context("when input="+input, func() {
					It("returns error "+resultErr.Error(), func() {
						_, err := eval(input)
						Expect(err).To(MatchError(resultErr))
					})
				})
			}
		})
	})
})
