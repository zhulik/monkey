package tokens

type TokenType string

const (
	// Can have literal.
	IDENTIFIER TokenType = "IDENTIFIER"
	INTEGER    TokenType = "INTEGER"
	STRING     TokenType = "STRING"

	// Literal is equal to the type itself.
	ASSIGN   TokenType = "="
	PLUS     TokenType = "+"
	MINUS    TokenType = "-"
	BANG     TokenType = "!"
	ASTERISK TokenType = "*"
	SLASH    TokenType = "/"

	GT  TokenType = ">"
	LT  TokenType = "<"
	GTE TokenType = ">="
	LTE TokenType = "<="

	EQ  TokenType = "=="
	NEQ TokenType = "!="

	COMMA     TokenType = ","
	SEMICOLON TokenType = ";"

	LPAREN TokenType = "("
	RPAREN TokenType = ")"
	LBRACE TokenType = "{"
	RBRACE TokenType = "}"

	FUNCTION TokenType = "fn"
	LET      TokenType = "let"
	TRUE     TokenType = "true"
	FALSE    TokenType = "false"
	IF       TokenType = "if" //nolint:varnamelen
	ELSE     TokenType = "else"
	RETURN   TokenType = "return"
	NIL      TokenType = "nil"
)

var keywords = map[TokenType]TokenType{ //nolint:gochecknoglobals
	FUNCTION: FUNCTION,
	LET:      LET,
	TRUE:     TRUE,
	FALSE:    FALSE,
	IF:       IF,
	ELSE:     ELSE,
	RETURN:   RETURN,
	NIL:      NIL,
}

type Token struct {
	Type    TokenType
	literal string
}

func New(tokenType TokenType, literal ...string) Token {
	lit := ""
	if len(literal) > 0 {
		lit = literal[0]
	}

	return Token{Type: tokenType, literal: lit}
}

func TypeFromLiteral(lit string) TokenType {
	if t, ok := keywords[TokenType(lit)]; ok {
		return t
	}

	return IDENTIFIER
}

func (t Token) Literal() string {
	if t.literal != "" {
		return t.literal
	}

	return string(t.Type)
}
