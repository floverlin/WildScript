package cmd

import (
	"fmt"
	"arc/cmd/cli"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize project",
	Long:  `Initialize project.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("creating arc project...")
		cli.InitProject()
		fmt.Println("file main.arc created")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
