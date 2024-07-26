package cmd

import (
	"vault-wars/app"
	"vault-wars/util"

	"github.com/spf13/cobra"
)

var nuke bool

var rmCmd = &cobra.Command{
	Use:   "rm",
	Short: "Uninstall a Vault helm chart on the cluster.",
	Run: func(cmd *cobra.Command, args []string) {
		if err := app.RemoveChart("luke", "default", nuke); err != nil {
			util.ExitError(err)
		}
	},
}

func init() {
	rmCmd.Flags().BoolVarP(&nuke, "nuke", "n", false, "Remove all Vault data in addition to uninstalling the chart")
	rootCmd.AddCommand(rmCmd)
}
