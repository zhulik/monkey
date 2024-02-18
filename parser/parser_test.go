package parser_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/zhulik/monkey/ast"
	"github.com/zhulik/monkey/lexer"
	"github.com/zhulik/monkey/parser"
	"github.com/zhulik/monkey/tokens"
)

func parseProgram(program string) (*ast.Program, error) {
	lex := lexer.New(program)

	par, err := parser.New(lex)
	if err != nil {
		return nil, err
	}

	return par.ParseProgram()
}

var _ = Describe("Parser", func() {
	Describe(".ParseProgram", func() {
		Context("when parsing let statements", func() {
			Context("when statements are valid", func() {
				input := `let foobar = 12345;
					let a = 1;
					let b = 2;`

				It("returns the let statements", func() {
					program, err := parseProgram(input)
					Expect(err).ToNot(HaveOccurred())
					Expect(program).NotTo(BeNil())

					statements := program.Statements

					lets := []struct {
						name  string
						value any
					}{
						{"foobar", 12345},
						{"a", 1},
						{"b", 2},
					}

					Expect(statements).To(HaveLen(len(lets)))
					for i, let := range lets {
						letStatement, ok := statements[i].(*ast.LetStatement)
						Expect(ok).To(BeTrue())

						Expect(letStatement.Token.Type).To(Equal(tokens.LET))
						Expect(letStatement.Name.Value).To(Equal(let.name))
						// Expect(letStatement.Value).To(Equal(12345))
					}
				})
			})

			Context("when the identitier is missed", func() {
				input := `let = 12345;`

				It("returns en error", func() {
					_, err := parseProgram(input)
					Expect(err).To(MatchError("invalid token. Expected: IDENTIFIER, found: ="))
				})
			})

			Context("when the assign sign is missed", func() {
				input := `let foobar 123;`

				It("returns en error", func() {
					_, err := parseProgram(input)
					Expect(err).To(MatchError("invalid token. Expected: =, found: INTEGER"))
				})
			})
		})

		Context("when parsing return statements", func() {
			Context("when statements are valid", func() {
				input := `return 123;
				return 234;
				return 456;`

				It("returns the return statements", func() {
					program, err := parseProgram(input)

					Expect(err).ToNot(HaveOccurred())
					Expect(program).NotTo(BeNil())

					statements := program.Statements

					returns := []struct{}{
						{},
						{},
						{},
					}

					Expect(statements).To(HaveLen(len(returns)))

					for i := range returns {
						letStatement, ok := statements[i].(*ast.ReturnStatement)
						Expect(ok).To(BeTrue())

						Expect(letStatement.Token.Type).To(Equal(tokens.RETURN))
						// Expect(letStatement.Name.Value).To(Equal(let.name))
						// Expect(letStatement.Value).To(Equal(12345))
					}
				})
			})

			Context("when statement is invalid", func() {
				// TODO: write me
			})
		})

		Context("when parsing identifier expressions", func() {
			Context("when expression is valid", func() {
				input := "foobar;"

				It("returns parsed expression", func() {
					program, err := parseProgram(input)
					Expect(err).ToNot(HaveOccurred())

					statements := program.Statements
					Expect(statements).To(HaveLen(1))

					stmt, ok := statements[0].(*ast.ExpressionStatement)
					Expect(ok).To(BeTrue())

					expr, ok := stmt.Expression.(*ast.IdentifierExpression)
					Expect(ok).To(BeTrue())

					Expect(expr.Value).To(Equal("foobar"))
					Expect(expr.TokenLiteral()).To(Equal("foobar"))
				})
			})

			Context("when expression is invalid", func() {
				// TODO: write me
			})
		})

		Context("when parsing integer expressions", func() {
			Context("when expression is valid", func() {
				input := "12345;"

				It("returns parsed expression", func() {
					program, err := parseProgram(input)
					Expect(err).ToNot(HaveOccurred())

					statements := program.Statements
					Expect(statements).To(HaveLen(1))

					stmt, ok := statements[0].(*ast.ExpressionStatement)
					Expect(ok).To(BeTrue())

					expr, ok := stmt.Expression.(*ast.IntegerExpression)
					Expect(ok).To(BeTrue())

					Expect(expr.Value).To(Equal(int64(12345)))
					Expect(expr.TokenLiteral()).To(Equal("12345"))
				})
			})

			Context("when expression is invalid", func() {
				// TODO: write me
			})
		})

		Context("when parsing prefix expressions", func() {
			Context("when expression is valid", func() {
				It("returns parsed expression", func() {
					cases := []struct {
						input    string
						operator string
						value    int64
					}{
						{"!5;", "!", 5},
						{"-15;", "-", 15},
					}

					for _, testCase := range cases {
						program, err := parseProgram(testCase.input)
						Expect(err).ToNot(HaveOccurred())

						statements := program.Statements
						Expect(statements).To(HaveLen(1))

						stmt, ok := statements[0].(*ast.ExpressionStatement)
						Expect(ok).To(BeTrue())

						expr, ok := stmt.Expression.(*ast.PrefixExpression)
						Expect(ok).To(BeTrue())

						Expect(expr.Operator).To(Equal(testCase.operator))

						intExpr, ok := expr.Right.(*ast.IntegerExpression)
						Expect(ok).To(BeTrue())

						Expect(intExpr.Value).To(Equal(testCase.value))
					}
				})
			})

			Context("when expression is invalid", func() {
				// TODO: write me
			})
		})

		Context("when parsing infix expressions", func() {
			Context("when expression is valid", func() {
				It("returns parsed expression", func() {
					cases := []struct {
						input      string
						leftValue  int64
						operator   string
						rightValue int64
					}{
						{"5+5;", 5, "+", 5},
						{"5-5;", 5, "-", 5},
						{"5*5;", 5, "*", 5},
						{"5/5;", 5, "/", 5},
						{"5>5;", 5, ">", 5},
						{"5<5;", 5, "<", 5},
						{"5>=5;", 5, ">=", 5},
						{"5<=5;", 5, "<=", 5},
						{"5==5;", 5, "==", 5},
						{"5!=5;", 5, "!=", 5},
					}

					for _, testCase := range cases {
						program, err := parseProgram(testCase.input)
						Expect(err).ToNot(HaveOccurred())

						statements := program.Statements
						Expect(statements).To(HaveLen(1))

						stmt, ok := statements[0].(*ast.ExpressionStatement)
						Expect(ok).To(BeTrue())

						expr, ok := stmt.Expression.(*ast.InfixExpression)
						Expect(ok).To(BeTrue())

						Expect(expr.Operator).To(Equal(testCase.operator))

						leftExpr, ok := expr.Right.(*ast.IntegerExpression)
						Expect(ok).To(BeTrue())

						rightExpr, ok := expr.Right.(*ast.IntegerExpression)
						Expect(ok).To(BeTrue())

						Expect(leftExpr.Value).To(Equal(testCase.leftValue))
						Expect(rightExpr.Value).To(Equal(testCase.rightValue))
					}
				})
			})

			Context("when expression is invalid", func() {
				// TODO: write me
			})
		})
	})
})
