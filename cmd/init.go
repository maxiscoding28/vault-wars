package cmd

import (
	"vault-wars/app"
	"vault-wars/util"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a deployed Vault cluster.",
	Run: func(cmd *cobra.Command, args []string) {
		if err := app.InitializeCluster("luke", "default"); err != nil {
			util.ExitError(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
