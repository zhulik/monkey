package parser_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/zhulik/monkey/lexer"
	"github.com/zhulik/monkey/parser"
)

var _ = Describe("Parser", func() {
	Describe(".ParseProgram", func() {
		Context("when program is valid", func() {
			cases := map[string]string{
				// Let statements.
				"let a = 5;":       "let a = 5;",
				"let foo = bar;":   "let foo = bar;",
				"let a = bar + 5;": "let a = (bar + 5);",

				// // return statements.
				"return 5;":       "return 5;",
				"return bar;":     "return bar;",
				"return bar + 5;": "return (bar + 5);",

				"nil": "nil",

				// Basic expressions.
				"foobar":         "foobar",
				"12345":          "12345",
				"!5":             "(!5)",
				"-15":            "(-15)",
				"5 + 5":          "(5 + 5)",
				"5 - 5":          "(5 - 5)",
				"5 * 5":          "(5 * 5)",
				"5 / 5":          "(5 / 5)",
				"5 > 5":          "(5 > 5)",
				"5 < 5":          "(5 < 5)",
				"5 >= 5":         "(5 >= 5)",
				"5 <= 5":         "(5 <= 5)",
				"5 == 5":         "(5 == 5)",
				"5 != 5":         "(5 != 5)",
				"true":           "true",
				"false":          "false",
				"2 > 3 == false": "((2 > 3) == false)",

				// Expressions and operator precedence.
				"-a * b":                     "((-a) * b)",
				"!-a":                        "(!(-a))",
				"a + b + c":                  "((a + b) + c)",
				"a + b - c":                  "((a + b) - c)",
				"a * b * c":                  "((a * b) * c)",
				"a * b / c":                  "((a * b) / c)",
				"a + b * c + d / e - f":      "(((a + (b * c)) + (d / e)) - f)",
				"3 + 4; -5 * 5":              "(3 + 4)((-5) * 5)",
				"5 > 4 == 3 < 4":             "((5 > 4) == (3 < 4))",
				"5 > 4 != 3 > 4":             "((5 > 4) != (3 > 4))",
				"3 + 4 * 5 == 3 * 1 + 4 * 5": "((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
				"3 + 4 * 5 >= 3 * 1 + 4 * 5": "((3 + (4 * 5)) >= ((3 * 1) + (4 * 5)))",
				"3 + 4 * 5 <= 3 * 1 + 4 * 5": "((3 + (4 * 5)) <= ((3 * 1) + (4 * 5)))",
				"1 + (2 + 3) + 4":            "((1 + (2 + 3)) + 4)",
				"(5 + 5) * 2":                "((5 + 5) * 2)",
				"2 / (5 + 5)":                "(2 / (5 + 5))",
				"-(5 + 5)":                   "(-(5 + 5))",
				"!(true == true)":            "(!(true == true))",
				"a + add(b * c) + d":         "((a + add((b * c))) + d)",
				"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))": "add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))",
				"add(a + b + c * d / f + g)":                "add((((a + b) + ((c * d) / f)) + g))",

				// If expressions.
				"if (x < y) { x }":            "if (x < y) { x }",
				"if (x < y) { }":              "if (x < y) { }",
				"if (x < y) { x } else { y }": "if (x < y) { x } else { y }",
				"if (x < y) { x } else { }":   "if (x < y) { x } else { }",

				// Functions.
				"fn(x, y) { x + y }": "fn(x, y) { (x + y) }",
				"fn() { 1 }":         "fn() { 1 }",
				"fn(x, y, z) { }":    "fn(x, y, z) { }",

				// Function calls.
				"foo()":        "foo()",
				"foo(1, 2, 3)": "foo(1, 2, 3)",
				`"foo bar"`:    `"foo bar"`,
			}

			for input, output := range cases {
				Context("when parsing "+input, func() {
					It("returns parsed "+output, func() {
						lex := lexer.New(input)

						par := parser.New(lex)

						program, err := par.ParseProgram()
						Expect(err).ToNot(HaveOccurred())
						Expect(program.String()).To(Equal(output))
					})
				})
			}
		})

		Context("when program is invalid", func() {
			// TODO: write me
		})
	})
})
