package wildscript

import (
	"fmt"
	"log"
	"os"
	"time"
	"wildscript/internal/lexer"
	"wildscript/internal/logger"
	"wildscript/internal/parser"
)

func Run() {
	args := os.Args
	if len(args) != 2 {
		log.Fatal("no file")
	}
	start := time.Now()
	input, err := os.ReadFile(args[1])
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

	fmt.Println(program)

	log.Printf(
		"[wild] program ends in %d us\n",
		time.Since(start).Microseconds(),
	)
}
