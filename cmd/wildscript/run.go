package wildscript

import (
	"log"
	"os"
	"time"
	"wildscript/internal/lexer"
	"wildscript/internal/logger"
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
	}

	log.Printf(
		"[wild] program ends in %d us\n",
		time.Since(start).Microseconds(),
	)
}
