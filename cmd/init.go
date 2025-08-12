package cmd

import (
	"fmt"
	"wildscript/cmd/cli"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initial",
	Long:  `Initial your project`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("creating wildscript project...")
		cli.InitProject()
		fmt.Println("file main.ws created")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
