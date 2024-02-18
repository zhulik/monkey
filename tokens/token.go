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
	LET      TokenType = "LET"
)

var keywords = map[string]TokenType{ //nolint:gochecknoglobals
	"fn":  FUNCTION,
	"let": LET,
}

type Token struct {
	Type    TokenType
	Literal string
}

func New(tokenType TokenType, literal ...string) Token {
	lit := ""
	if len(literal) > 0 {
		lit = literal[0]
	}

	return Token{Type: tokenType, Literal: lit}
}

func TypeFromLiteral(lit string) TokenType {
	if t, ok := keywords[lit]; ok {
		return t
	}

	return IDENTIFIER
}
