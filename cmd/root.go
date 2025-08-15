package cmd

import (
	"arc/cmd/interpreter"
	"arc/internal/settings"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "arc",
	Short: "Arc interpreter",
	Long:  `Arc is a minimalist programming language.`,
	Args:  cobra.ArbitraryArgs,

	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("no file")
			cmd.Usage()
			return
		}
		
		var file string
		ext := filepath.Ext(args[0])
		if ext == "" {
			file = args[0] + ".arc"
		} else if ext != ".arc" {
			fmt.Printf(
				"wrong file type %s\n",
				ext,
			)
			cmd.Usage()
			return
		} else {
			file = args[0]
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
