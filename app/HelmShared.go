package app

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

func isReleaseNotDeployed(releaseName string) error {
	cmd := exec.Command("helm", "list", "-o", "json")
	out, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("error running the helm list command: %v", err)
	}

	var releases []HelmRelease
	err = json.Unmarshal(out, &releases)
	if err != nil {
		return fmt.Errorf("error parsing JSON output: %v", err)
	}

	for _, release := range releases {
		if release.Name == releaseName {
			return fmt.Errorf("release named \"%s\" is already deployed", releaseName)
		}
	}

	// Release is not found, return nil
	return nil
}

func isReleaseDeployed(releaseName string) error {
	cmd := exec.Command("helm", "list", "-o", "json")
	out, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("error running the helm list command: %v", err)
	}

	var releases []HelmRelease
	err = json.Unmarshal(out, &releases)
	if err != nil {
		return fmt.Errorf("error parsing JSON output: %v", err)
	}

	for _, release := range releases {
		if release.Name == releaseName {
			return nil
		}
	}

	// Release is not found, return an error
	return fmt.Errorf("no release deployed with the following name: \"%s\"", releaseName)
}
