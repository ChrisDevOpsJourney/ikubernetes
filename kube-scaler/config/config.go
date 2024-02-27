package config

import "github.com/kelseyhightower/envconfig"

type Config struct {
	UseInClusterClient bool   `envconfig:"USE_IN_CLUSTER_CLIENT"`
	Namespace          string `envconfig:"NAMESPACE"`
	ResourceKind       string `envconfig:"RESOURCE_KIND"`
	Replicas           int32  `envconfig:"REPLICAS"`
}

func LoadConfig() (*Config, error) {
	var config Config

	err := envconfig.Process("", &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
