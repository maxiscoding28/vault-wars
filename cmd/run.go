package cmd

import (
	"vault-wars/app"
	"vault-wars/util"

	"github.com/spf13/cobra"
)

var deployCmd = &cobra.Command{
	Use:   "run",
	Short: "Install a Vault helm chart on the cluster.",
	Run: func(cmd *cobra.Command, args []string) {
		if err := app.DeployChart("luke", "default"); err != nil {
			util.ExitError(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(deployCmd)
}
