package cmd

import (
	"fmt"
	"wildscript/cmd/cli"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize",
	Long:  `Initialize new project.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("creating sigillum project...")
		cli.InitProject()
		fmt.Println("file main.sil created")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
