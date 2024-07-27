package app

import (
	"fmt"
	"os"
	"path/filepath"
	"vault-wars/util"
)

func InitializeCluster(releaseName string, namespace string) error {
	if err := isReleaseDeployed(releaseName); err != nil {
		return err
	}

	k8Client, k8Config, err := createKubernetesClient()
	if err != nil {
		return err
	}

	if err := AllPodsRunning(k8Client, releaseName, namespace); err != nil {
		return err
	}

	maxRetries := 3
	err = EnsureNoPodsInitialized(k8Client, k8Config, releaseName, namespace, maxRetries)
	if err != nil {
		return err
	}
	util.LogInfo(fmt.Sprintf("No pods already initialized after %d checks.", maxRetries))
	util.LogInfo(fmt.Sprintf("Attempting to initialize node %s", util.InitNodeName(releaseName)))

	command := "vault operator init -key-shares=1 -key-threshold=1 -format=json"
	out, err := execOnPod(k8Client, k8Config, util.InitNodeName(releaseName), namespace, "vault", command, false)
	if err != nil {
		return err
	}
	util.LogInfo(fmt.Sprintf("Success! node %s initialized", util.InitNodeName(releaseName)))
	util.LogInfo(out)

	fileName := fmt.Sprintf("%s-cluster-init.json", releaseName)
	if err := util.WriteFile(fileName, out); err != nil {
		return err
	}
	workingDirectory, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting working directory: %s", err)
	}
	filePath := filepath.Join(workingDirectory, fileName)
	util.LogInfo(fmt.Sprintf("Unseal keys and root token saved at:\n\t%s", filePath))
	return nil
}
