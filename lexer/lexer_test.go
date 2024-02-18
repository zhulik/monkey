package lexer_test

import (
	"io"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/zhulik/monkey/lexer"
	"github.com/zhulik/monkey/tokens"
)

var _ = Describe("Lexer", func() {
	Describe(".NextToken", func() {
		input := `;`

		Context("when there is a next token", func() {
			It("returns the next token", func() {
				lex := lexer.New(input)
				token, err := lex.NextToken()
				Expect(err).ToNot(HaveOccurred())
				Expect(token).To(Equal(tokens.New(tokens.SEMICOLON)))
			})
		})

		Context("when lexer is at the end of input", func() {
			It("returns io.EOF", func() {
				lex := lexer.New(input)

				token, err := lex.NextToken()
				Expect(err).ToNot(HaveOccurred())
				Expect(token).To(Equal(tokens.New(tokens.SEMICOLON)))

				_, err = lex.NextToken()
				Expect(err).To(MatchError(io.EOF))

				_, err = lex.NextToken()
				Expect(err).To(MatchError(io.EOF))
			})
		})
	})
	Describe(".Tokens", func() {
		Context("when all tokens are correct", func() {
			input := `=+(){},;
let five = 5;
let ten = 10;

let add = fn(x, y) {
	x + y;
};

let result = add(five, ten);

!-/*5;
5 < 10 > 5;
if (5 < 10) {
  return true;
} else {
  return false;
}
10 == 10;
10 != 9;
10 >= 10;
9 <= 10;`
			lex := lexer.New(input)

			result := []tokens.Token{
				tokens.New(tokens.ASSIGN),
				tokens.New(tokens.PLUS),
				tokens.New(tokens.LPAREN),
				tokens.New(tokens.RPAREN),
				tokens.New(tokens.LBRACE),
				tokens.New(tokens.RBRACE),
				tokens.New(tokens.COMMA),
				tokens.New(tokens.SEMICOLON),

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

				tokens.New(tokens.INTEGER, "10"),
				tokens.New(tokens.EQ),
				tokens.New(tokens.INTEGER, "10"),
				tokens.New(tokens.SEMICOLON),

				tokens.New(tokens.INTEGER, "10"),
				tokens.New(tokens.NEQ),
				tokens.New(tokens.INTEGER, "9"),
				tokens.New(tokens.SEMICOLON),

				tokens.New(tokens.INTEGER, "10"),
				tokens.New(tokens.GTE),
				tokens.New(tokens.INTEGER, "10"),
				tokens.New(tokens.SEMICOLON),

				tokens.New(tokens.INTEGER, "9"),
				tokens.New(tokens.LTE),
				tokens.New(tokens.INTEGER, "10"),
				tokens.New(tokens.SEMICOLON),
			}

			It("returns the next token", func() {
				tokens, err := lex.Tokens()
				Expect(err).ToNot(HaveOccurred())
				Expect(tokens).To(Equal(result))
			})
		})

		Context("when there is an illegal character", func() {
			lex := lexer.New("$#")

			It("returns an error", func() {
				_, err := lex.Tokens()
				Expect(err).To(MatchError(lexer.ErrIllegalCharacter))
			})
		})

		Context("when parsing an empty string", func() {
			lex := lexer.New("")

			It("returns an empty array", func() {
				tokens, err := lex.Tokens()
				Expect(tokens).To(BeEmpty())
				Expect(err).ToNot(HaveOccurred())
			})
		})
	})
})
