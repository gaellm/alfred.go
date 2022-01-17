package conf

import (
	"strings"

	"github.com/spf13/viper"
)

const (
	//Mock files directory configuration key name.
	MOCKS_DIR = "core.mocks-dir"
)

// Struct where all config keys are stored.
type Config struct {
	//Mocks' folder.
	Core CoreConfig `mapstructure:"core"`
}

// Struct where all core config keys are stored.
type CoreConfig struct {
	MocksDir string `mapstructure:"mocks-dir"`
}

//Configure the Viper instance.
func configureViper(v *viper.Viper) (*viper.Viper, error) {

	v.SetConfigName("config")
	v.SetConfigType("toml")
	v.AddConfigPath("configs")

	err := v.ReadInConfig() // Find and read the config file
	if err != nil {
		return v, err
	}

	// SetEnvKeyReplacer sets the strings.Replacer on the viper object
	// Useful for mapping an environmental variable to a key that does
	// not match it
	// needed to allow viper to replace . in nested variables by _ when looking at environment variable
	replacer := strings.NewReplacer(".", "_", "-", "_")
	v.SetEnvKeyReplacer(replacer)

	// Environnement variables superseed the values retrieved from config file
	v.AutomaticEnv()

	return v, err
}

//Use Viper to get configuration from files and environment variables.
func buildConfiguration() (Config, error) {

	//Initiate a Viper instance.
	v := viper.New()

	configuration := Config{}

	v, err := configureViper(v)
	if err != nil {
		return Config{}, err
	}

	//Unmarshall loaded config into our struct
	err = v.Unmarshal(&configuration)
	if err != nil {
		return Config{}, err
	}

	return configuration, nil
}
