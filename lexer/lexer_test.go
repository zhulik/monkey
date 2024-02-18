package lexer_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/zhulik/monkey/lexer"
	"github.com/zhulik/monkey/tokens"
)

var _ = Describe("Lexer", func() {
	Describe(".NextToken", func() {
		Context("when lexing simple tokens", func() {
			input := "=+(){},;"
			lex := lexer.New(input)

			tests := []tokens.Token{
				tokens.New(tokens.ASSIGN),
				tokens.New(tokens.PLUS),
				tokens.New(tokens.LPAREN),
				tokens.New(tokens.RPAREN),
				tokens.New(tokens.LBRACE),
				tokens.New(tokens.RBRACE),
				tokens.New(tokens.COMMA),
				tokens.New(tokens.SEMICOLON),
				tokens.New(tokens.EOF),
			}

			It("returns the next token", func() {
				for _, testCase := range tests {
					Expect(lex.NextToken()).To(Equal(testCase))
				}
			})
		})

		Context("when lexing more complicated tokens", func() {
			input := `let five = 5;
let ten = 10;

let add = fn(x, y) {
	x + y;
};

let result = add(five, ten);`

			lex := lexer.New(input)

			tests := []tokens.Token{
				tokens.New(tokens.LET),
				tokens.New(tokens.IDENTIFIER, "five"),
				tokens.New(tokens.ASSIGN),
				tokens.New(tokens.INTEGER, "5"),
				tokens.New(tokens.SEMICOLON),

				tokens.New(tokens.LET),
				tokens.New(tokens.IDENTIFIER, "ten"),
				tokens.New(tokens.ASSIGN),
				tokens.New(tokens.INTEGER, "10"),
				tokens.New(tokens.SEMICOLON),

				tokens.New(tokens.LET),
				tokens.New(tokens.IDENTIFIER, "add"),
				tokens.New(tokens.ASSIGN),
				tokens.New(tokens.FUNCTION),
				tokens.New(tokens.LPAREN),
				tokens.New(tokens.IDENTIFIER, "x"),
				tokens.New(tokens.COMMA),
				tokens.New(tokens.IDENTIFIER, "y"),
				tokens.New(tokens.RPAREN),
				tokens.New(tokens.LBRACE),
				tokens.New(tokens.IDENTIFIER, "x"),
				tokens.New(tokens.PLUS),
				tokens.New(tokens.IDENTIFIER, "y"),
				tokens.New(tokens.SEMICOLON),
				tokens.New(tokens.RBRACE),
				tokens.New(tokens.SEMICOLON),

				tokens.New(tokens.LET),
				tokens.New(tokens.IDENTIFIER, "result"),
				tokens.New(tokens.ASSIGN),
				tokens.New(tokens.IDENTIFIER, "add"),
				tokens.New(tokens.LPAREN),
				tokens.New(tokens.IDENTIFIER, "five"),
				tokens.New(tokens.COMMA),
				tokens.New(tokens.IDENTIFIER, "ten"),
				tokens.New(tokens.RPAREN),
				tokens.New(tokens.SEMICOLON),
				tokens.New(tokens.EOF),
			}

			It("returns the next token", func() {
				for _, testCase := range tests {
					Expect(lex.NextToken()).To(Equal(testCase))
				}
			})
		})

		Context("when lexing some other tokens", func() {
			input := `!-/*5;
5 < 10 > 5;`

			lex := lexer.New(input)

			tests := []tokens.Token{
				tokens.New(tokens.BANG),
				tokens.New(tokens.MINUS),
				tokens.New(tokens.SLASH),
				tokens.New(tokens.ASTERISK),
				tokens.New(tokens.INTEGER, "5"),
				tokens.New(tokens.SEMICOLON),

				tokens.New(tokens.INTEGER, "5"),
				tokens.New(tokens.LT),
				tokens.New(tokens.INTEGER, "10"),
				tokens.New(tokens.GT),
				tokens.New(tokens.INTEGER, "5"),
				tokens.New(tokens.SEMICOLON),
				tokens.New(tokens.EOF),
			}

			It("returns the next token", func() {
				for _, testCase := range tests {
					Expect(lex.NextToken()).To(Equal(testCase))
				}
			})
		})

		Context("when lexing some other tokens", func() {
			input := `if (5 < 10) {
  return true;
} else {
  return false;
}`

			lex := lexer.New(input)

			tests := []tokens.Token{
				tokens.New(tokens.IF),
				tokens.New(tokens.LPAREN),
				tokens.New(tokens.INTEGER, "5"),
				tokens.New(tokens.LT),
				tokens.New(tokens.INTEGER, "10"),
				tokens.New(tokens.RPAREN),
				tokens.New(tokens.LBRACE),
				tokens.New(tokens.RETURN),
				tokens.New(tokens.TRUE),
				tokens.New(tokens.SEMICOLON),
				tokens.New(tokens.RBRACE),
				tokens.New(tokens.ELSE),
				tokens.New(tokens.LBRACE),
				tokens.New(tokens.RETURN),
				tokens.New(tokens.FALSE),
				tokens.New(tokens.SEMICOLON),
				tokens.New(tokens.RBRACE),
				tokens.New(tokens.EOF),
			}

			It("returns the next token", func() {
				for _, testCase := range tests {
					Expect(lex.NextToken()).To(Equal(testCase))
				}
			})
		})
	})
})
