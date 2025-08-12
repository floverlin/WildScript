package cli

import "os"

func InitProject() {
	fileName := "main.ws"
	program := "print(\"hello, world!\")"
	os.WriteFile(fileName, []byte(program), 0644)
}
