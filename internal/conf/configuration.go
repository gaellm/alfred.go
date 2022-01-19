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
