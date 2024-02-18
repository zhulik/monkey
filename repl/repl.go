package repl

import (
	"errors"
	"fmt"
	"io"

	"github.com/chzyer/readline"
	"github.com/k0kubun/pp"
	"github.com/zhulik/monkey/lexer"
)

func Start() error {
	pp.Printf("Monkey repl.\n")

	rln, err := readline.New(">> ")
	if err != nil {
		return fmt.Errorf("readline init error: %w", err)
	}

	defer rln.Close()

	for {
		line, rErr := rln.Readline()
		if rErr != nil {
			if errors.Is(rErr, io.EOF) || errors.Is(rErr, readline.ErrInterrupt) {
				return nil
			}

			return fmt.Errorf("readline error: %w", rErr)
		}

		lex := lexer.New(line)
		// TODO: use rangefunc when golangci-lint implemetns proper support
		// for _, token := range lex.IterateTokens() {
		// 	fmt.Printf("%+v\n", token)
		// }

		for token, lErr := lex.NextToken(); !errors.Is(lErr, io.EOF); token, lErr = lex.NextToken() {
			if lErr != nil {
				pp.Print("Lexing error: %w", lErr)

				break
			}

			pp.Printf("%+v\n", token)
		}
	}
}
