package cmd

import (
	"fmt"
	"log"
	"os"
	"time"
	"wildscript/internal/evaluator"
	"wildscript/internal/lexer"
	"wildscript/internal/logger"
	"wildscript/internal/parser"
	"wildscript/internal/settings"
	"wildscript/pkg"
)

func runFile(fileName string) {
	start := time.Now()

	gs := settings.Global
	input, err := os.ReadFile(fileName)
	if err != nil {
		log.Fatal("read file error: ", err)
	}

	l := lexer.New(input)
	c := lexer.NewCollector(l)

	if len(c.Illegals()) != 0 {
		for _, illegal := range c.Illegals() {
			logger.Log(
				illegal.Line,
				illegal.Column,
				"[lexer] illegal token: %s",
				illegal.Literal,
			)
		}
		os.Exit(1)
	}

	p := parser.New(c)
	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		for _, err := range p.Errors() {
			fmt.Printf(
				"[parser] error: %s",
				err,
			)
		}
		os.Exit(1)
	}

	if gs.Debug {
		fmt.Print(
			pkg.Cover(
				program.String(),
				"program",
				"-",
				20,
			),
		)
	}

	e := evaluator.New()

	defer wrapPanic()

	if !gs.Debug {
		e.Eval(program)
		return
	}

	for idx, stmt := range program.Statements {
		obj := e.Eval(stmt)
		fmt.Printf("%d >> %s\n", idx+1, obj.Inspect())
	}

	fmt.Printf(
		"[wild] program ends in %d us\n",
		time.Since(start).Microseconds(),
	)
}

func wrapPanic() {
	if p := recover(); p != nil {
		fmt.Printf("[wild] runtime error: %s", p)
		os.Exit(1)
	}
}
