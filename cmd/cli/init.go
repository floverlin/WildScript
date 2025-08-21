package cli

import (
	_ "embed"
	"os"
)

//go:embed assets/main.sil
var file []byte

func InitProject() {
	fileName := "main.sil"
	os.WriteFile(fileName, file, 0644)
}
