package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"vault-wars/util"
)

func UnsealCluster(releaseName string, namespace string) error {
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

	podInitialized, err := IsPodInitialized(k8Client, k8Config, util.InitNodeName(releaseName), namespace)
	if err != nil {
		return err
	}
	if !podInitialized {
		return fmt.Errorf("%s not initialized. Can't unseal", util.InitNodeName(releaseName))
	}

	data, err := util.ReadFile(fmt.Sprintf("%s-cluster-init.json", releaseName))
	if err != nil {
		return err
	}

	var vaultInit VaultInitOutput
	if err := json.Unmarshal(data, &vaultInit); err != nil {
		return fmt.Errorf("error unmarshalling JSON: %v", err)
	}

	unsealKey := vaultInit.UnsealKeysHex[0]
	rootToken := vaultInit.RootToken

	command := fmt.Sprintf("vault operator unseal %s", unsealKey)
	out, err := execOnPod(k8Client, k8Config, util.InitNodeName(releaseName), namespace, "vault", command, false)
	if err != nil {
		return fmt.Errorf("error executing initial vault unseal command on pod %s: %v", util.InitNodeName(releaseName), err)
	}
	fmt.Printf("Initial unseal command executed on pod %s\n%v", util.InitNodeName(releaseName), string(out))

	totalPods, initializedPods, err := CountInitializedPods(k8Client, k8Config, releaseName, namespace)
	if err != nil {
		return err
	}

	if totalPods != initializedPods {
		return errors.New("not all pods are initialized. Retry in a few seconds")
	}

	if err := UnsealPods(k8Client, k8Config, releaseName, namespace, util.InitNodeName(releaseName), unsealKey); err != nil {
		return err
	}

	util.LogInfo("Success. Cluster unsealed")
	util.LogInfo(fmt.Sprintf("Root token: %s", rootToken))

	return nil
}
