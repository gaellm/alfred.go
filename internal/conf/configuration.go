package conf

import (
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

	return configuration, nil
}
