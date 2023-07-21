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

package server

import (
	"alfred/internal/conf"
	"alfred/internal/function"
	"alfred/internal/helper"
	"alfred/internal/log"
	"alfred/internal/mock"
	"alfred/internal/tracing"
	"alfred/pkg/metrics"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"
)

type Key string

func pathHelperMiddleware(mockCollection mock.MockCollection) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			for _, m := range mockCollection.Mocks {
				if m.HasRegexUrl() {
					if m.Request.RegexUrl.Match([]byte(r.URL.Path)) {

						//var values Values
						values := make(map[string]string)

						//add the helper name:value in the context
						hs := m.GetPathRegexHelpers()

						//populate
						for _, h := range hs {

							index, err := strconv.Atoi(h.Target)
							if err != nil {
								log.Error(r.Context(), "Failed to transform helper target into integer type", err)
							}

							if err != nil {
								fmt.Println("Error during conversion")
								return
							}

							v := m.Request.RegexUrl.FindSubmatch([]byte(r.URL.Path))

							//add to values
							values[h.String] = string(v[index])
						}

						//add to original path
						values["originalPath"] = r.URL.Path

						//update url to match the mock handler
						r.URL.Path = m.Request.UrlTransformed

						ctx := context.WithValue(r.Context(), helper.PathHelperKey("pathHelperValues"), values)
						next.ServeHTTP(w, r.WithContext(ctx))
						return
					}
				}
			}

			// call next handler
			next.ServeHTTP(w, r)

		}
		return http.HandlerFunc(fn)
	}
}

func routerMiddleware(next http.Handler) http.Handler {

	fn := func(w http.ResponseWriter, r *http.Request) {

		r.URL.Path = "/" + r.Method + r.URL.Path

		// call next handler
		next.ServeHTTP(w, r)

	}
	return http.HandlerFunc(fn)

}

func removeFirstFolder(path string) string {
	components := strings.SplitN(path, "/", 3)
	if len(components) >= 3 {
		return "/" + components[2]
	}
	return path
}

// Middleware for logging each request
func logRequestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Serve the request
		next.ServeHTTP(w, r)

		var path string
		values := r.Context().Value(helper.PathHelperKey("pathHelperValues"))
		if values != nil {

			path = values.(map[string]string)["originalPath"]

		} else {
			path = r.URL.Path
		}

		// Log the request
		msg := " " + r.Method + " " + removeFirstFolder(path) + " " + r.RemoteAddr + " " + fmt.Sprint(time.Since(start))

		log.Info(r.Context(), msg)

		if log.GetLevel() == "debug" {

			reqBodyBytes, _ := io.ReadAll(r.Body)

			log.Debug(r.Context(), msg,
				zap.String("request-body", string(reqBodyBytes)))
		}
	})
}

// Build the service
func BuildServer(conf *conf.Config, asyncRunningJobsCount *sync.WaitGroup, mockCollection mock.MockCollection) (*http.Server, error) {
	//Build all endpoints handler
	handler, err := BuildHandler(conf, asyncRunningJobsCount, mockCollection)
	if err != nil {
		return nil, err
	}

	//Router
	handler = routerMiddleware(handler)

	//logger
	handler = logRequestMiddleware(handler)

	//tracing
	handler = tracing.AddTracingMiddlware(handler)

	//
	if mockCollection.HasRegexUrlMock() {
		log.Debug(context.Background(), "Regex URL detected")
		middlewarePathHelper := pathHelperMiddleware(mockCollection)
		handler = middlewarePathHelper(handler)
	}

	//Associate the handler to a server (-> contains listening interface(s))
	//Here, the server is listening on ALL interfaces and binding on 'conf.Port' port
	return &http.Server{
		Handler: handler,
		Addr:    fmt.Sprintf("%s:%s", conf.Alfred.Core.Listen.Ip, conf.Alfred.Core.Listen.Port),
	}, nil
}

// Create Handler with endpoints to serve: used from service or from tests
func BuildHandler(conf *conf.Config, asyncRunningJobsCount *sync.WaitGroup, mocks mock.MockCollection) (http.Handler, error) {
	//log
	log.Info(context.Background(), "Starting")
	defer log.LogPanic()

	// Create controller
	mux := http.NewServeMux()

	if conf.Alfred.Prometheus.Enable {

		prometheusConfig := metrics.MetricsConfig{
			MetricPath:     conf.Alfred.Prometheus.Path,
			MetricPort:     conf.Alfred.Prometheus.Listen.Port,
			MetricIp:       conf.Alfred.Prometheus.Listen.Ip,
			HttpServerIp:   conf.Alfred.Core.Listen.Ip,
			HttpServerPort: conf.Alfred.Core.Listen.Port,
			SlowTime:       conf.Alfred.Prometheus.SlowTimeSeconds,
		}
		prometheusConfig.SanitizeConfiguration()
		//metrics.CreateMetricEngine(controller, prometheusConfig)
		//metrics
		metrics.AddMetrics(mux, prometheusConfig)
		log.Info(context.Background(), "Prometheus exporter started to serve on host "+prometheusConfig.MetricIp+" and is listening at port "+fmt.Sprint(prometheusConfig.MetricPort)+" with '"+prometheusConfig.MetricPath+"' path")
	}

	// Global delay applied to mocks
	var alfredGlobalDelay time.Duration

	{

		//Add Routes
		{
			mux.HandleFunc("/POST"+"/logger", ChangingLoggingLevelRuntime)

			mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

				mockList, _ := json.MarshalIndent(mocks.GetMockInfoList(), "", "   ")
				_, err := w.Write([]byte("Hello Sir ! I take care ot the following mocks:\n" + string(mockList) + "\n\n(Alfred)"))
				if err != nil {
					log.Error(r.Context(), "failed to write", err)
				}

			})

			mux.HandleFunc("/PATCH"+"/alfred", func(w http.ResponseWriter, r *http.Request) {
				PatchMock(w, r, mocks)
			})

			mux.HandleFunc("/POST"+"/alfred/delay", func(w http.ResponseWriter, r *http.Request) {
				DelayMocks(&alfredGlobalDelay, w, r)
			})

			//Load JS functions
			functionCollection, err := function.CreateFunctionCollectionFromFolder(conf.Alfred.Core.FunctionsDir)
			if err != nil {
				log.Debug(context.Background(), "function files loader error: "+err.Error())
			}

			// Create mocks routes
			AddMocksRoutes(mux, mocks, functionCollection, &alfredGlobalDelay)
		}
	}

	return mux, nil
}

// Serve will bind the port(s) and launch serve in a separated goroutine
func Serve(main_ctx context.Context, conf *conf.Config, server *http.Server) {

	//Bind
	listener, err := net.Listen("tcp", server.Addr)
	if err != nil {
		log.Error(main_ctx, "error Server Binding", err)
		panic("error Server Binding on " + server.Addr)
	}

	//Log that's the bind is ok
	fmt.Println("{ \"alfred-speaking\" : \"Started to serve on host " + conf.Alfred.Core.Listen.Ip + " and listening at port " + conf.Alfred.Core.Listen.Port + ", with " + main_ctx.Value(Key("mocksNb")).(string) + " mocks, Sir.\"}")
	log.Info(main_ctx, "alfred started to serve on host "+conf.Alfred.Core.Listen.Ip+" and is listening at port "+conf.Alfred.Core.Listen.Port)

	//Start to serve
	go func() {

		//tracing
		cleanup, err := tracing.Init(main_ctx, tracing.OtelConfig{
			ServiceName:           conf.Alfred.Name,
			ServiceNamespace:      conf.Alfred.Namespace,
			DeploymentEnvironment: conf.Alfred.Environment,
			ExporterInsecure:      conf.Alfred.Tracing.Insecure,
			TracesSampler:         conf.Alfred.Tracing.Sampler,
			TracesSamplerArg:      conf.Alfred.Tracing.SamplerArgs,
			ExporterOtlpEndpoint:  conf.Alfred.Tracing.OtlpEndpoint,
		})
		if err != nil {
			log.Error(main_ctx, "server panic", errors.New("error during preparing tracer..."+err.Error()))
			panic("error during preparing tracer..." + err.Error())
		}
		defer func() {
			err := cleanup(main_ctx)
			if err != nil {
				log.Error(main_ctx, "error during tracing cleanup", err)
			}
		}()

		err = server.Serve(listener)
		if err != nil {
			log.Error(main_ctx, "error at server start", err)
		}
	}()
}

// Return necessaries to stop serve
func Stop(ctx context.Context, server *http.Server, asyncRunningJobsCount *sync.WaitGroup) {
	log.Info(ctx, "Server is stopping")

	//Let's some few seconds to shutdown gracefully
	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(5*time.Second))
	defer cancel()

	//Shutdown the http server
	if err := server.Shutdown(ctx); err != nil {
		log.Error(ctx, "Error While stopping Server: ", err)
	}

	// Wait that all async jobs are done (timeboxed)
	waitAsyncJobsTimeout(ctx, asyncRunningJobsCount)

	//Log again...
	log.Info(ctx, "Stopped Server")
}

// waitAsyncJobsTimeout is synchronized on 2 things :
// - a channel to be aware that's an async job is finished
// - a timeout context to timebox this 'waiting'
func waitAsyncJobsTimeout(ctx context.Context, wg *sync.WaitGroup) bool {
	c := make(chan struct{})
	go func() {
		defer close(c)
		wg.Wait()
	}()
	select {
	case <-c:
		return false // wait() was done
	case <-ctx.Done():
		return true // deadline was reached
	}
}

// requestRecover retreives a panic if it happened during request treatment
func requestRecover(w http.ResponseWriter, req *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Error(req.Context(), "Catched Panic", errors.New(fmt.Sprint(r)))
		}
	}()
}
