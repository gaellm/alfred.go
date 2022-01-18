package conf

import (
	"strings"

	"github.com/imdario/mergo"
)

const (
	DEFAULT_MOCKS_DIR = "user-files/mocks/"
)

var DefaultConfig = Config{
	Core: CoreConfig{
		MocksDir: DEFAULT_MOCKS_DIR,
	},
}

func mergeConfigs(dst *Config, src Config) error {

	if err := mergo.Merge(dst, src); err != nil {
		return err
	}

	return nil
}

//Create and init the configuration.
func GetConfiguration() (Config, error) {

	configuration, err := buildConfiguration()
	if err != nil {
		return Config{}, err
	}

	if mergeConfigs(&configuration, DefaultConfig); err != nil {
		return Config{}, err
	}

	return sanitizeConfiguration(configuration), nil
}

//Check and correct configuration values
func sanitizeConfiguration(configuration Config) Config {

	if !strings.HasSuffix(configuration.Core.MocksDir, "/") {

		configuration.Core.MocksDir += "/"
	}

	return configuration
}
