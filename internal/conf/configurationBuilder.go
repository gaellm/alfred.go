/*
 * Copyright The Alfred.go Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package conf

import (
	"strings"

	"github.com/spf13/viper"
)

const (
	//Mock files directory configuration key name.
	MOCKS_DIR_KEY = "alfred.core.mocks-dir"

	//Component name configuration key name.
	NAME_KEY = "alfred.name"

	//Component version configuration key name.
	VERSION_KEY = "alfred.version"

	//Listening interface key name.
	LISTEN_INTERFACE = "alfred.core.listen.ip"

	//Listening port key name.
	LISTEN_PORT = "alfred.core.listen.port"
)

// Struct where all config keys are stored.
type Config struct {
	// Whole configuration.
	Alfred AlfredConfig `mapstructure:"alfred"`
}

// Struct where all config keys are stored.
type AlfredConfig struct {
	Name    string `mapstructure:"name"`
	Version string `mapstructure:"version"`
	//Core configuration.
	Core CoreConfig `mapstructure:"core"`
}

type ListenConfig struct {
	Ip   string `mapstructure:"ip"`
	Port string `mapstructure:"port"`
}

// Struct where all core config keys are stored.
type CoreConfig struct {
	MocksDir string       `mapstructure:"mocks-dir"`
	Listen   ListenConfig `mapstructure:"listen"`
}

//Configure the Viper instance.
func configureViper(v *viper.Viper) (*viper.Viper, error) {

	v.SetConfigName("config")
	v.SetConfigType("json")
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

	// Set default is needed by viper to init the key,
	// if not the env variable will not be applied if
	// configuration key is commented.
	v.SetDefault(NAME_KEY, "")
	v.SetDefault(MOCKS_DIR_KEY, "")
	v.SetDefault(VERSION_KEY, "")
	//v.SetDefault(LISTEN_INTERFACE, "")
	//v.SetDefault(LISTEN_PORT, "")

	//Unmarshall loaded config into our struct
	err = v.Unmarshal(&configuration)
	if err != nil {
		return Config{}, err
	}

	return configuration, nil
}
