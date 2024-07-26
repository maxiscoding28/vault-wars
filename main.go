package main

// func main() {
// 	if err := app.VerifyDependencies(); err != nil {
// 		util.ExitError(err)
// 	}

// 	// vwrs deploy luke
// 	// optional flags --namespace default --seal-type shamir
// 	if err := app.DeployChart("luke", "default"); err != nil {
// 		util.ExitError(err)
// 	}

// 	// vwrs initialize luke
// 	// optional flag --namespace
// 	if err := app.InitializeCluster("luke", "default"); err != nil {
// 		util.ExitError(err)
// 	}

// 	// vwrs unseal luke
// 	// optional flag --namespace
// 	if err := app.UnsealCluster("luke", "default"); err != nil {
// 		util.ExitError(err)
// 	}

// 	// vwrs destroy luke
// 	// optional flag --namespace
// 	// if err := app.DestroyCluster("luke", "default"); err != nil {

// 	// }

// 	// Minikube tunnel
// 	// To get VAULT_ADDR
// }

import (
	"vault-wars/cmd"
)

func main() {
	cmd.Execute()
}
