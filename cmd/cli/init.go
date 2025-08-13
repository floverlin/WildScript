package cli

import (
	_ "embed"
	"os"
)

//go:embed assets/main.ws
var file []byte

func InitProject() {
	fileName := "main.ws"
	os.WriteFile(fileName, file, 0644)
}
