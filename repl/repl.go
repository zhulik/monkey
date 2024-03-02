package repl

import (
	"errors"
	"fmt"
	"io"

	"github.com/chzyer/readline"
	"github.com/k0kubun/pp"
	"github.com/zhulik/monkey/evaluator"
	obj "github.com/zhulik/monkey/evaluator/object"
	"github.com/zhulik/monkey/lexer"
	"github.com/zhulik/monkey/parser"
)

func Start() error {
	pp.Printf("Monkey repl.\n")

	eval := evaluator.New()

	rln, err := readline.New(">> ")
	if err != nil {
		return fmt.Errorf("readline init error: %w", err)
	}

	defer rln.Close()

	environment := obj.NewEnv()

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
			fmt.Printf("Parsing error: %s\n", pErr.Error()) //nolint:forbidigo

			continue
		}

		result, eErr := eval.Eval(program, environment)
		if eErr != nil {
			fmt.Printf("Evaluation error: %s\n", eErr.Error()) //nolint:forbidigo

			continue
		}

		pp.Println(result.Inspect())
	}
}
