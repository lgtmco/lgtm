package model

type DeploymentConfigs map[string]DeploymentConfig

type DeploymentConfig struct {
	Tasks          []string    `toml:"tasks"`
	Environment    *string     `toml:"env"`
}

type DeploymentInfo struct {
	Ref         string
	Task        string
	Environment string
}
