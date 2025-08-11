package cmd

import (
	"fmt"
	"os"
	"wildscript/internal/settings"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "wild",
	Short: "WildScript interpreter",
	Long:  `WildScript - GO, GO WILD!`,

	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("no file")
			cmd.Usage()
			return
		}
		runFile(args[0])
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolVarP(
		&settings.Global.Debug,
		"debug",
		"d",
		false,
		"enable debug mode",
	)
}
