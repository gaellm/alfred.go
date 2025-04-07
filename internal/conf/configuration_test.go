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
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestDefaultConfig(t *testing.T) {
	// Verify default configuration values
	assert.Equal(t, DEFAULT_NAME, DefaultConfig.Alfred.Name)
	assert.Equal(t, DEFAULT_VERSION, DefaultConfig.Alfred.Version)
	assert.Equal(t, DEFAULT_NAMESPACE, DefaultConfig.Alfred.Namespace)
	assert.Equal(t, DEFAULT_ENVIRONMENT, DefaultConfig.Alfred.Environment)
	assert.Equal(t, DEFAULT_MOCKS_DIR, DefaultConfig.Alfred.Core.MocksDir)
	assert.Equal(t, DEFAULT_FUNCTIONS_DIR, DefaultConfig.Alfred.Core.FunctionsDir)
	assert.Equal(t, DEFAULT_BODIES_DIR, DefaultConfig.Alfred.Core.BodiesDir)
	assert.Equal(t, DEFAULT_LISTEN_INTERFACE, DefaultConfig.Alfred.Core.Listen.Ip)
	assert.Equal(t, DEFAULT_LISTEN_PORT, DefaultConfig.Alfred.Core.Listen.Port)
	assert.Equal(t, DEFAULT_TLS_ENABLED, DefaultConfig.Alfred.Core.Listen.TlsEnabled)
	assert.Equal(t, DEFAULT_TLS_CERT_PATH, DefaultConfig.Alfred.Core.Listen.TlsCertPath)
	assert.Equal(t, DEFAULT_TLS_KEY_PATH, DefaultConfig.Alfred.Core.Listen.TlsKeyPath)
	assert.Equal(t, DEFAULT_LOG_LEVEL, DefaultConfig.Alfred.LogLevel)
	assert.Equal(t, DEFAULT_PROMETHEUS_ENABLE, DefaultConfig.Alfred.Prometheus.Enable)
	assert.Equal(t, DEFAULT_PROMETHEUS_PATH, DefaultConfig.Alfred.Prometheus.Path)
	assert.Equal(t, DEFAULT_PROMETHEUS_SLOW_TIME_SECONDS, DefaultConfig.Alfred.Prometheus.SlowTimeSeconds)
	assert.Equal(t, DEFAULT_TRACING_OTLP_ENDPOINT, DefaultConfig.Alfred.Tracing.OtlpEndpoint)
	assert.Equal(t, DEFAULT_TRACING_INSECURE, DefaultConfig.Alfred.Tracing.Insecure)
	assert.Equal(t, DEFAULT_TRACING_SAMPLER, DefaultConfig.Alfred.Tracing.Sampler)
	assert.Equal(t, DEFAULT_TRACING_SAMPLER_ARGS, DefaultConfig.Alfred.Tracing.SamplerArgs)
}

func TestMergeConfigs(t *testing.T) {
	// Create a source configuration with some custom values
	srcConfig := Config{
		Alfred: AlfredConfig{
			Name:    "custom-name",
			Version: "2.0",
			Core: CoreConfig{
				MocksDir: "custom-mocks/",
			},
		},
	}

	// Create a destination configuration with default values
	defautCfg := DefaultConfig

	// Merge the configurations
	err := mergeConfigs(&srcConfig, defautCfg)
	assert.NoError(t, err)

	// Verify that the merged configuration contains the custom values
	assert.Equal(t, "info", srcConfig.Alfred.LogLevel)

	// Verify that other values remain unchanged
	assert.Equal(t, "custom-name", srcConfig.Alfred.Name)

}

func TestSanitizeConfiguration(t *testing.T) {
	// Create a configuration with a missing trailing slash in MocksDir
	config := Config{
		Alfred: AlfredConfig{
			Core: CoreConfig{
				MocksDir: "user-files/mocks",
			},
		},
	}

	// Sanitize the configuration
	sanitizedConfig := sanitizeConfiguration(config)

	// Verify that the trailing slash is added
	assert.Equal(t, "user-files/mocks/", sanitizedConfig.Alfred.Core.MocksDir)
}

func TestBuildConfiguration(t *testing.T) {
	// Create a temporary configuration file
	configContent := `
	{
		"alfred": {
			"name": "test-name",
			"core": {
				"mocks-dir": "test-mocks/"
			}
		}
	}`
	configFile := "config.json"
	configFolder := "configs"
	configPath := filepath.Join(configFolder, configFile)
	err := os.Mkdir(configFolder, 0755) // Create the configs folder
	assert.NoError(t, err)
	err = os.WriteFile(configPath, []byte(configContent), 0644)
	assert.NoError(t, err)
	defer os.RemoveAll(configFolder)

	// Build the configuration
	config, err := buildConfiguration()
	assert.NoError(t, err)

	// Verify that the configuration contains the values from the file
	assert.Equal(t, "test-name", config.Alfred.Name)
	assert.Equal(t, "test-mocks/", config.Alfred.Core.MocksDir)
}

func TestConfigureViper(t *testing.T) {
	// Set an environment variable
	os.Setenv("ALFRED_CORE_MOCKS_DIR", "env-mocks/")
	defer os.Unsetenv("ALFRED_CORE_MOCKS_DIR")

	// Create a temporary configuration file
	configContent := `
	{
		"alfred": {
			"name": "test-name",
			"core": {
				"mocks-dir": "test-mocks/"
			}
		}
	}`
	configFile := "config.json"
	configFolder := "configs"
	configPath := filepath.Join(configFolder, configFile)
	err := os.Mkdir(configFolder, 0755) // Create the configs folder
	assert.NoError(t, err)
	err = os.WriteFile(configPath, []byte(configContent), 0644)
	assert.NoError(t, err)
	defer os.RemoveAll(configFolder)

	// Configure Viper
	v := viper.New()
	v, err = configureViper(v)
	assert.NoError(t, err)

	// Verify that the environment variable is applied
	assert.Equal(t, "env-mocks/", v.GetString(MOCKS_DIR_KEY))
}

func TestGetConfiguration(t *testing.T) {
	// Create a temporary configuration file

	configContent := `
		{
			"alfred": {
				"name": "test-name",
				"core": {
					"mocks-dir": "test-mocks/"
				}
			}
		}`
	configFile := "config.json"
	configFolder := "configs"
	configPath := filepath.Join(configFolder, configFile)
	err := os.Mkdir(configFolder, 0755) // Create the configs folder
	assert.NoError(t, err)
	err = os.WriteFile(configPath, []byte(configContent), 0644)
	assert.NoError(t, err)
	defer os.RemoveAll(configFolder)

	// Get the configuration
	config, err := GetConfiguration()
	assert.NoError(t, err)

	// Verify that the configuration contains the values from the file
	assert.Equal(t, "test-name", config.Alfred.Name)
	assert.Equal(t, "test-mocks/", config.Alfred.Core.MocksDir) // Ensure sanitizeConfiguration is applied
}
