package cli

import (
	_ "embed"
	"os"
)

//go:embed assets/main.wild
var file []byte

func InitProject() {
	fileName := "main.wild"
	os.WriteFile(fileName, file, 0644)
}
