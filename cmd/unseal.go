package cmd

import (
	"vault-wars/app"
	"vault-wars/util"

	"github.com/spf13/cobra"
)

var unseal = &cobra.Command{
	Use:   "unseal",
	Short: "Install a Vault helm chart on the cluster.",
	Run: func(cmd *cobra.Command, args []string) {
		if err := app.UnsealCluster("luke", "default"); err != nil {
			util.ExitError(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(unseal)
}
