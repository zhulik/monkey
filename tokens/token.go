package tokens

type TokenType string

const (
	ILLEGAL TokenType = "ILLEGAL"
	EOF     TokenType = "EOF"

	IDENTIFIER TokenType = "IDENTIFIER"
	INTEGER    TokenType = "INTEGER"

	ASSIGN TokenType = "="
	PLUS   TokenType = "+"

	COMMA     TokenType = ","
	SEMICOLON TokenType = ";"

	LPAREN TokenType = "("
	RPAREN TokenType = ")"
	LBRACE TokenType = "{"
	RBRACE TokenType = "}"

	FUNCTION TokenType = "FUNCTION"
	LET      TokenType = "LEN"
)

type Token struct {
	Type    TokenType
	Literal string
}
