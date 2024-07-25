package app

import (
	"vault-wars/util"
	"fmt"
	"encoding/json"
)

func isReleaseDeployed(releaseName string) error {
	out, err := util.ExecCommand("helm", "list", "-o", "json")
	if err != nil {
		return fmt.Errorf("error running the helm list command: %v", out)
	}
	var releases []HelmRelease
	err = json.Unmarshal(out, &releases)
	if err != nil {
		return fmt.Errorf("error parsing JSON output: %v", err)
	}

	for _, release := range releases {
		if release.Name == releaseName {
			return fmt.Errorf("a release named \"%s\" is already deployed", releaseName)
		}
	}

	return nil
}