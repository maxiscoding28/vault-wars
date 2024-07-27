package app

import (
	"errors"
	"fmt"
	"vault-wars/util"
)

func DeployChart(releaseName string, namespace string) error {
	if err := isReleaseNotDeployed(releaseName); err != nil {
		return err
	}
	if err := installChart(releaseName, namespace); err != nil {
		return err
	}

	return nil
}

func installChart(releaseName string, namespace string) error {
	if !isValidReleaseName(releaseName) {
		return fmt.Errorf("release name %s is not valid\nValid names are: luke, leia, anakin", releaseName)
	}
	filename := fmt.Sprintf("%s.yaml", releaseName)
	if err := util.WriteFile(filename, YamlValuesMap[releaseName]); err != nil {
		return err
	}
	if out, err := util.ExecCommand("helm", "install", releaseName, "hashicorp/vault", "--values", filename, "--namespace", namespace); err != nil {
		return errors.New(string(out))
	}
	util.LogInfo("Helm chart installed successfully.")
	return nil
}

func isValidReleaseName(releaseName string) bool {
	_, exists := YamlValuesMap[releaseName]
	return exists
}
