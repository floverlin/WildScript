package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"wildscript/cmd/interpreter"
	"wildscript/internal/settings"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "sigil",
	Short: "Sigillum interpreter",
	Long:  `Sigillum - witchcraft language`,
	Args:  cobra.ArbitraryArgs,

	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("no file")
			cmd.Usage()
			return
		}

		var file string
		if ext := filepath.Ext(args[0]); ext == "" {
			file = args[0] + ".sil" // TODO CONST
		}

		interpreter.RunFile(file)
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolVar(
		&settings.Global.Debug,
		"debug",
		false,
		"enable debug mode",
	)
	rootCmd.Flags().BoolVar(
		&settings.Global.Tokens,
		"tokens",
		false,
		"show lexer tokens",
	)
}
