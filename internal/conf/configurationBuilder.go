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

	//Functions files directory configuration key name.
	FUNCTIONS_DIR_KEY = "alfred.core.functions-dir"

	//Body files directory configuration key name.
	BODIES_DIR_KEY = "alfred.core.body-files-dir"

	//Component name configuration key name.
	NAME_KEY = "alfred.name"

	//Component version configuration key name.
	VERSION_KEY = "alfred.version"

	NAMESPACE_KEY   = "alfred.namespace"
	ENVIRONMENT_KEY = "alfred.environment"

	//Listening interface key name.
	LISTEN_INTERFACE = "alfred.core.listen.ip"

	//Listening port key name.
	LISTEN_PORT          = "alfred.core.listen.port"
	LISTEN_TLS_ENABLE    = "alfred.core.listen.enable-tls"
	LISTEN_TLS_CERT_PATH = "alfred.core.listen.tls-cert-path"
	LISTEN_TLS_KEY_PATH  = "alfred.core.listen.tls-key-path"

	//Debug key
	LOG_LEVEL_KEY = "alfred.log-level"

	//Prometheus
	PROMETHEUS_ENABLE_KEY            = "alfred.prometheus.enable"
	PROMETHEUS_PATH_KEY              = "alfred.prometheus.path"
	PROMETHEUS_SLOW_TIME_SECONDS_KEY = "alfred.prometheus.slow-time-seconds"
	PROMETHEUS_LISTEN_PORT_KEY       = "alfred.prometheus.listen.port"
	PROMETHEUS_LISTEN_IP_KEY         = "alfred.prometheus.listen.ip"

	//Tracing
	TRACING_OTLP_ENDPOINT_KEY = "alfred.tracing.otlp-endpoint"
	TRACING_INSECURE_KEY      = "alfred.tracing.insecure"
	TRACING_SAMPLER_KEY       = "alfred.tracing.sampler"
	TRACING_SAMPLER_ARGS_KEY  = "alfred.tracing.sampler-args"
)

// Struct where all config keys are stored.
type Config struct {
	// Whole configuration.
	Alfred AlfredConfig `mapstructure:"alfred"`
}

// Struct where all config keys are stored.
type AlfredConfig struct {
	Name        string `mapstructure:"name"`
	Version     string `mapstructure:"version"`
	Namespace   string `mapstructure:"namespace"`
	Environment string `mapstructure:"environment"`
	LogLevel    string `mapstructure:"log-level"`
	//Core configuration.
	Core       CoreConfig       `mapstructure:"core"`
	Prometheus PrometheusConfig `mapstructure:"prometheus"`
	Tracing    TracingConfig    `mapstructure:"tracing"`
}

type ListenConfig struct {
	Ip          string `mapstructure:"ip"`
	Port        string `mapstructure:"port"`
	TlsEnabled  bool   `mapstructure:"enable-tls"`
	TlsCertPath string `mapstructure:"tls-cert-path"`
	TlsKeyPath  string `mapstructure:"tls-key-path"`
}

// Struct where all core config keys are stored.
type CoreConfig struct {
	MocksDir     string       `mapstructure:"mocks-dir"`
	FunctionsDir string       `mapstructure:"functions-dir"`
	BodiesDir    string       `mapstructure:"body-files-dir"`
	Listen       ListenConfig `mapstructure:"listen"`
}

type PrometheusConfig struct {
	Enable          bool         `mapstructure:"enable"`
	Path            string       `mapstructure:"path"`
	SlowTimeSeconds int32        `mapstructure:"slow-time-seconds"`
	Listen          ListenConfig `mapstructure:"listen"`
}

type TracingConfig struct {
	OtlpEndpoint string `mapstructure:"otlp-endpoint"`
	Insecure     bool   `mapstructure:"insecure"`
	Sampler      string `mapstructure:"sampler"`
	SamplerArgs  string `mapstructure:"sampler-args"`
}

// Configure the Viper instance.
func configureViper(v *viper.Viper) (*viper.Viper, error) {

	v.SetConfigName("config")
	v.SetConfigType("json")
	v.AddConfigPath("configs")

	// SetEnvKeyReplacer sets the strings.Replacer on the viper object
	// Useful for mapping an environmental variable to a key that does
	// not match it
	// needed to allow viper to replace . in nested variables by _ when looking at environment variable
	replacer := strings.NewReplacer(".", "_", "-", "_")
	v.SetEnvKeyReplacer(replacer)

	// Environnement variables superseed the values retrieved from config file
	v.AutomaticEnv()

	err := v.ReadInConfig() // Find and read the config file
	if err != nil {
		return v, err
	}

	return v, err
}

// Use Viper to get configuration from files and environment variables.
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
	v.SetDefault(FUNCTIONS_DIR_KEY, "")
	v.SetDefault(BODIES_DIR_KEY, "")
	v.SetDefault(VERSION_KEY, "")
	v.SetDefault(NAMESPACE_KEY, "")
	v.SetDefault(ENVIRONMENT_KEY, "")
	v.SetDefault(LISTEN_INTERFACE, "")
	v.SetDefault(LISTEN_PORT, "")
	v.SetDefault(LISTEN_TLS_ENABLE, "")
	v.SetDefault(LISTEN_TLS_CERT_PATH, "")
	v.SetDefault(LISTEN_TLS_KEY_PATH, "")
	v.SetDefault(LOG_LEVEL_KEY, "")
	v.SetDefault(PROMETHEUS_ENABLE_KEY, "")
	v.SetDefault(PROMETHEUS_PATH_KEY, "")
	v.SetDefault(PROMETHEUS_SLOW_TIME_SECONDS_KEY, "")
	v.SetDefault(PROMETHEUS_LISTEN_PORT_KEY, "")
	v.SetDefault(PROMETHEUS_LISTEN_IP_KEY, "")
	v.SetDefault(TRACING_OTLP_ENDPOINT_KEY, "")
	v.SetDefault(TRACING_INSECURE_KEY, "")
	v.SetDefault(TRACING_SAMPLER_KEY, "")
	v.SetDefault(TRACING_SAMPLER_ARGS_KEY, "")

	//Unmarshall loaded config into our struct
	err = v.Unmarshal(&configuration)
	if err != nil {
		return Config{}, err
	}

	return configuration, nil
}
