package cmd

import (
	"wildscript/cmd/cli"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initial",
	Long:  `Initial your project`,
	Run: func(cmd *cobra.Command, args []string) {
		cli.InitProject()
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
