package app

type Repo struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type HelmRelease struct {
	Name string `json:"name"`
}

type VaultInitOutput struct {
	UnsealKeysHex []string `json:"unseal_keys_hex"`
	RootToken     string   `json:"root_token"`
}

type VaultStatus struct {
	Initialized bool `json:"initialized"`
}
