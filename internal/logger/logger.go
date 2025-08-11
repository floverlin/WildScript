package logger

import "fmt"

func Log(line, col int, str string, args ...any) {
	init := fmt.Sprintf(str, args...)
	fmt.Printf("%s at line: %d col: %d\n", init, line, col)
}

func Slog(line, col int, str string, args ...any) string {
	init := fmt.Sprintf(str, args...)
	return fmt.Sprintf("%s at line: %d col: %d\n", init, line, col)
}
