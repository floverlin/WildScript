package cmd

import (
	"fmt"
	"os"
	"wildscript/cmd/interpreter"
	"wildscript/internal/settings"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "wild",
	Short: "WildScript interpreter",
	Long:  `WildScript - GO, GO WILD!`,
	Args:  cobra.ArbitraryArgs,

	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("no file")
			cmd.Usage()
			return
		}
		interpreter.RunFile(args[0])
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
