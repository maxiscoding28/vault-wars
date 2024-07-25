package app

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"vault-wars/util"
)

func InitializeCluster(releaseName string, namespace string) error {
	if err := isReleaseDeployed(releaseName); err == nil {
		return fmt.Errorf("a release named %s is not deployed", releaseName)
	}

	k8Client, k8Config, err := createKubernetesClient()
	if err != nil {
		return err
	}

	if err := AllPodsRunning(k8Client, releaseName, namespace); err != nil {
		return err
	}

	_, initializedPods, err := CountInitializedPods(k8Client, k8Config, releaseName, namespace)
	if err != nil {
		return err
	}

	if initializedPods > 0 {
		return errors.New("unable to initialize node. Nodes `are already initialized")
	}

	command := []string{"/bin/sh", "-c", "vault operator init -key-shares=1 -key-threshold=1 -format=json"}
	out, err := execOnPod(k8Client, k8Config, util.InitNodeName(releaseName), namespace, "vault", command, false)
	if err != nil {
		return err
	}
	util.LogInfo(fmt.Sprintf("Success! %s-cluster initialized", releaseName))
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
	util.LogInfo(fmt.Sprintf("Unseal keys and root token are save at:%s", filePath))
	return nil
}
