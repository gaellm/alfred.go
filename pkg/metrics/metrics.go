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
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/penglongli/gin-metrics/ginmetrics"
)

type MetricsConfig struct {
	MetricPath     string
	MetricPort     string
	MetricIp       string
	HttpServerIp   string
	HttpServerPort string
	SlowTime       int32
	Logger         gin.HandlerFunc
}

func CreateMetricEngine(e *gin.Engine, conf MetricsConfig) {

	conf.SanitizeConfiguration()

	appRouter := e
	metricRouter := appRouter
	m := ginmetrics.GetMonitor()

	dedicatedServer := conf.HttpServerIp != conf.MetricIp || conf.HttpServerPort != conf.MetricPort

	if dedicatedServer {
		metricRouter = gin.New()
		metricRouter.Use(conf.Logger)
		// use metric middleware without expose metric path
		m.UseWithoutExposingEndpoint(appRouter)
	} else {
		m.Use(metricRouter)
	}

	// set metric path expose to metric router
	m.SetMetricPath(conf.MetricPath)
	// +optional set slow time, default 5s
	m.SetSlowTime(conf.SlowTime)
	// +optional set request duration, default {0.1, 0.3, 1.2, 5, 10}
	// used to p95, p99
	//m.SetDuration([]float64{0.1, 0.3, 1.2, 5, 10})
	m.Expose(metricRouter)

	if dedicatedServer {
		go func() {
			_ = metricRouter.Run(fmt.Sprintf("%s:%s", conf.MetricIp, conf.MetricPort))
		}()
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
