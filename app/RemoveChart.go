package app

import (
	"fmt"
	"vault-wars/util"
)

func RemoveChart(releaseName string, namespace string, nuke bool) error {
	if err := isReleaseDeployed(releaseName); err == nil {
		util.LogError(fmt.Sprintf("no release named %s is deployed", releaseName))
	} else {
		if err := uninstallChart(releaseName, namespace); err != nil {
			return err
		}
	}

	if nuke {
		util.LogWarn("Nuke activated. Deleting PVCs")
		out, err := util.ExecCommand("kubectl", "delete", "pvc", "-l", fmt.Sprintf("app.kubernetes.io/instance=%s", releaseName))
		if err != nil {
			return err
		}
		util.LogInfo(fmt.Sprintf("deleting PVCs... \n%s", string(out)))
	}

	return nil
}

func uninstallChart(releaseName string, namespace string) error {
	if !isValidReleaseName(releaseName) {
		util.LogWarn(fmt.Sprintf("Release name %s is not valid.\nValid names are: luke, leia, anakin", releaseName))
		return fmt.Errorf("release name %s is not valid", releaseName)
	}
	filename := fmt.Sprintf("%s.yaml", releaseName)
	if err := util.WriteFile(filename, YamlValuesMap[releaseName]); err != nil {
		return err
	}
	if out, err := util.ExecCommand("helm", "uninstall", releaseName, "--namespace", namespace); err != nil {
		return fmt.Errorf("\n%s", string(out))
	}
	util.LogInfo("Helm chart uninstalled successfully.")
	return nil
}
