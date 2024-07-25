package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"vault-wars/util"
)

func VerifyDependencies() error {
	if !isBinaryInstalled("helm") {
		util.ExitError(errors.New("this project requires the helm binary accessible in your path"))
	}
	if !isBinaryInstalled("minikube") {
		util.ExitError(errors.New("this project requires the minikube binary accessible in your path"))
	}
	if !isMinikubeRunning() {
		util.ExitError(errors.New("minikube is not running\nTry running minikube start"))
	}
	if !helmRepoAdded("hashicorp") {
		if err := addHelmRepo("hashicorp"); err != nil {
			util.ExitError(err)
		}
	}
	return nil
}

func isBinaryInstalled(binaryName string) bool {
	if _, err := exec.LookPath(binaryName); err != nil {
		util.LogError(fmt.Sprintf("Unable to find %s binary, error: %d", binaryName, err.Error()))
		return false
	}
	util.LogInfo(fmt.Sprintf("%s is installed.", binaryName))
	return true
}

func isMinikubeRunning() bool {
	if out, err := util.ExecCommand("minikube", "status"); err != nil {
		util.LogError(fmt.Sprintf("minikube is not running - %s", string(out)))
		return false
	}
	util.LogInfo("minikube is running")
	return true
}

func helmRepoAdded(repoName string) bool {
	out, err := util.ExecCommand("helm", "repo", "list", "-o=json")
	if err != nil {
		util.LogError(fmt.Sprintf("Error running the helm repo list command - %s", string(out)))
		return false
	}
	var repos []Repo
	if err := json.Unmarshal(out, &repos); err != nil {
		util.LogError(fmt.Sprintf("Error unmarshalling json from helm repo list command - %s", err.Error()))
		return false
	}
	exists := false
	for _, repo := range repos {
		if repo.Name == repoName {
			exists = true
			break
		}
	}
	if err != nil {
		util.LogWarn(fmt.Sprintf("%s repo doesn't exist.", err.Error()))
		return false
	}

	return exists
}

func addHelmRepo(repoName string) error {
	util.LogWarn(fmt.Sprintf("%s repo doesn't exist. Attempting to add...", repoName))

	if out, err := util.ExecCommand("helm", "repo", "add", "hashicorp", "https://helm.releases.hashicorp.com"); err != nil {
		return errors.New(string(out))
	}
	return nil
}
