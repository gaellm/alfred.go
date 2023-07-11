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

package metrics

import (
	"alfred/internal/log"
	"context"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type MetricsConfig struct {
	MetricPath     string
	MetricPort     string
	MetricIp       string
	HttpServerIp   string
	HttpServerPort string
	SlowTime       int32
}

func AddMetrics(mux *http.ServeMux, conf MetricsConfig) {

	var metricsMux *http.ServeMux

	conf.SanitizeConfiguration()

	// Create non-global registry.
	reg := prometheus.NewRegistry()

	// Add go runtime metrics and process collectors.
	reg.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
	)

	promHandler := promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg})

	//create server
	dedicatedServer := conf.HttpServerIp != conf.MetricIp || conf.HttpServerPort != conf.MetricPort

	if dedicatedServer {
		metricsMux = http.NewServeMux()

		metricsMux.Handle(conf.MetricPath, promHandler)

		// Create a new HTTP server
		server := &http.Server{
			Addr:    conf.MetricIp + ":" + conf.MetricPort, // Specify the port to listen on
			Handler: metricsMux,                            // Use the ServeMux as the handler
		}

		go func() {
			err := server.ListenAndServe()
			if err != nil {
				log.Error(context.Background(), "serve metrics error", err)
			}
		}()

	} else {
		metricsMux = mux
		metricsMux.Handle(conf.MetricPath, promHandler)
	}
}

func (c *MetricsConfig) SanitizeConfiguration() {

	if c.MetricIp == "" {
		c.MetricIp = c.HttpServerIp
	}

	if c.MetricPort == "" {
		c.MetricPort = c.HttpServerPort
	}
}
