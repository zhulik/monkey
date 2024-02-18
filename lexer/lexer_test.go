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
				{Type: tokens.ASSIGN, Literal: "="},
				{Type: tokens.PLUS, Literal: "+"},
				{Type: tokens.LPAREN, Literal: "("},
				{Type: tokens.RPAREN, Literal: ")"},
				{Type: tokens.LBRACE, Literal: "{"},
				{Type: tokens.RBRACE, Literal: "}"},
				{Type: tokens.COMMA, Literal: ","},
				{Type: tokens.SEMICOLON, Literal: ";"},
				{Type: tokens.EOF},
			}

			It("returns the next token", func() {
				for _, testCase := range tests {
					Expect(lex.NextToken()).To(Equal(testCase))
				}
			})
		})
	})
})
