package interpreter

import (
	"fmt"
	"log"
	"os"
	"time"
	"wildscript/internal/environment"
	"wildscript/internal/evaluator"
	"wildscript/internal/lexer"
	"wildscript/internal/logger"
	"wildscript/internal/parser"
	"wildscript/internal/settings"
	"wildscript/pkg"

	"github.com/fatih/color"
)

func RunFile(fileName string) {
	start := time.Now()

	gs := settings.Global
	input, err := os.ReadFile(fileName)
	if err != nil {
		log.Fatal("read file error: ", err)
	}

	l := lexer.New(input)
	c := lexer.NewCollector(l)

	if gs.Debug && gs.Tokens {
		length := len(c.Tokens())
		for idx, token := range c.Tokens() {
			fmt.Print(token)
			if idx != length-1 {
				if (idx+1)%4 != 0 {
					fmt.Print(
						color.RedString(
							" | ",
						),
					)
				} else {
					fmt.Println()
				}

			}
		}
		fmt.Println()
	}

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

	defer wrapPanic()

	program := p.ParseProgram()

	if gs.Debug {
		fmt.Print(
			pkg.Cover(
				program.String(),
				"program",
				"-",
				32,
			),
		)
	}

	e := evaluator.New(nil)

	if !gs.Debug {
		e.Eval(program)
		return
	}

	var result environment.Object
	for idx, stmt := range program.Statements {
		obj := e.Eval(stmt)
		result = obj
		fmt.Printf("%d >> %s\n", idx+1, obj.Inspect())
	}
	fmt.Printf(
		"%s >>> %s\n",
		color.RedString("[program result]"),
		result.Inspect(),
	)

	fmt.Printf(
		"[sigil] program ends in %d us\n",
		time.Since(start).Microseconds(),
	)
}

func wrapPanic() {
	if p := recover(); p != nil {
		fmt.Printf("%s\n", p)
		os.Exit(1)
	}
}
