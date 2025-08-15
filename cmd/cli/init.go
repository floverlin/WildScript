package cli

import (
	_ "embed"
	"os"
)

//go:embed assets/main.arc
var file []byte

func InitProject() {
	fileName := "main.arc"
	os.WriteFile(fileName, file, 0644)
}
