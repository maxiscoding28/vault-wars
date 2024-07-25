package app

func DestroyCluster(releaseName string, namespace string) error {
	if err := isReleaseDeployed(releaseName); err != nil {
		return err
	}

    // helm uninstall deploy
    // check for error
    return nil
}