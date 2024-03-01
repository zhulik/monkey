package repl

import (
	"errors"
	"fmt"
	"io"

	"github.com/chzyer/readline"
	"github.com/k0kubun/pp"
	"github.com/zhulik/monkey/lexer"
	"github.com/zhulik/monkey/parser"
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

		parser := parser.New(lex)

		program, pErr := parser.ParseProgram()
		if pErr != nil {
			fmt.Printf("Parsing error: %s\n", pErr.Error())

			continue
		}

		pp.Printf("%+v\n", program.String())
	}
}
