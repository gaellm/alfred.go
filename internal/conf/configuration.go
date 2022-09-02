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
	"fmt"
	"strings"

	"github.com/imdario/mergo"
)

const (
	DEFAULT_NAME                         = "alfred-mock"
	DEFAULT_VERSION                      = "1.0"
	DEFAULT_NAMESPACE                    = "default"
	DEFAULT_ENVIRONMENT                  = "all"
	DEFAULT_MOCKS_DIR                    = "user-files/mocks/"
	DEFAULT_FUNCTIONS_DIR                = "user-files/functions/"
	DEFAULT_BODIES_DIR                   = "user-files/body-files/"
	DEFAULT_LISTEN_INTERFACE             = "0.0.0.0"
	DEFAULT_LISTEN_PORT                  = "8080"
	DEFAULT_LOG_LEVEL                    = "info"
	DEFAULT_PROMETHEUS_ENABLE            = false
	DEFAULT_PROMETHEUS_PATH              = "/metrics"
	DEFAULT_PROMETHEUS_SLOW_TIME_SECONDS = 10
	DEFAULT_PROMETHEUS_LISTEN_PORT       = ""
	DEFAULT_PROMETHEUS_LISTEN_IP         = ""
	DEFAULT_TRACING_OTLP_ENDPOINT        = ""
	DEFAULT_TRACING_INSECURE             = true
	DEFAULT_TRACING_SAMPLER              = "parentbased_traceidratio"
	DEFAULT_TRACING_SAMPLER_ARGS         = "1.0"
)

var DefaultConfig = Config{
	AlfredConfig{
		Name:        DEFAULT_NAME,
		Version:     DEFAULT_VERSION,
		Namespace:   DEFAULT_NAMESPACE,
		Environment: DEFAULT_ENVIRONMENT,
		LogLevel:    DEFAULT_LOG_LEVEL,
		Core: CoreConfig{
			MocksDir:     DEFAULT_MOCKS_DIR,
			FunctionsDir: DEFAULT_FUNCTIONS_DIR,
			BodiesDir:    DEFAULT_BODIES_DIR,
			Listen: ListenConfig{
				Ip:   DEFAULT_LISTEN_INTERFACE,
				Port: DEFAULT_LISTEN_PORT,
			},
		},
		Prometheus: PrometheusConfig{
			Enable:          DEFAULT_PROMETHEUS_ENABLE,
			Path:            DEFAULT_PROMETHEUS_PATH,
			SlowTimeSeconds: DEFAULT_PROMETHEUS_SLOW_TIME_SECONDS,
			Listen: ListenConfig{
				Ip:   DEFAULT_PROMETHEUS_LISTEN_IP,
				Port: DEFAULT_PROMETHEUS_LISTEN_PORT,
			},
		},
		Tracing: TracingConfig{
			OtlpEndpoint: DEFAULT_TRACING_OTLP_ENDPOINT,
			Insecure:     DEFAULT_TRACING_INSECURE,
			Sampler:      DEFAULT_TRACING_SAMPLER,
			SamplerArgs:  DEFAULT_TRACING_SAMPLER_ARGS,
		},
	},
}

func mergeConfigs(dst *Config, src Config) error {

	if err := mergo.Merge(dst, src); err != nil {
		return err
	}

	return nil
}

// Create and init the configuration.
func GetConfiguration() (Config, error) {

	configuration, err := buildConfiguration()
	if err != nil {
		fmt.Println("{ \"alfred-speaking\" : \"" + strings.Replace(err.Error(), "\"", "\\\"", -1) + ", gona use default configuration, Sir.\" }")
	}

	err = mergeConfigs(&configuration, DefaultConfig)
	if err != nil {
		return Config{}, err
	}

	return sanitizeConfiguration(configuration), nil
}

// Check and correct configuration values
func sanitizeConfiguration(configuration Config) Config {

	if !strings.HasSuffix(configuration.Alfred.Core.MocksDir, "/") {

		configuration.Alfred.Core.MocksDir += "/"
	}

	return configuration
}
