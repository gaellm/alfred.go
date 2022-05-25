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
	"alfred/internal/log"
	"alfred/internal/mock"
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

//Build the service
func BuildServer(conf *conf.Config, asyncRunningJobsCount *sync.WaitGroup, mocks mock.MockCollection) (*http.Server, error) {
	//Build all endpoints handler
	handler, err := BuildHandler(conf, asyncRunningJobsCount, mocks)
	if err != nil {
		return nil, err
	}

	//Associate the handler to a server (-> contains listening interface(s))
	//Here, the server is listening on ALL interfaces and binding on 'conf.Port' port
	return &http.Server{
		Handler: handler,
		Addr:    fmt.Sprintf("%s:%s", conf.Alfred.Core.Listen.Ip, conf.Alfred.Core.Listen.Port),
	}, nil
}

//Provide an handler to log access
func GinLogger() gin.HandlerFunc {

	return func(c *gin.Context) {

		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		method := c.Request.Method
		ip := c.ClientIP()
		proto := c.Request.Proto
		status := c.Writer.Status()
		c.Next()
		cost := time.Since(start)

		message := fmt.Sprintf("%s %s %s?%s %d",
			method,
			proto,
			path,
			query,
			status,
		)

		log.Info(context.Background(), message,
			zap.String("ip", ip),
			zap.String("status", fmt.Sprintf("%d", status)),
			zap.String("duration", cost.String()),
		)
	}
}

//Create Handler with endpoints to serve: used from service or from tests
func BuildHandler(conf *conf.Config, asyncRunningJobsCount *sync.WaitGroup, mocks mock.MockCollection) (http.Handler, error) {
	//log
	log.Info(context.Background(), "Starting")
	defer log.LogPanic()

	// Create controller
	gin.SetMode(gin.ReleaseMode)
	controller := gin.New()
	controller.Use(GinLogger())

	{
		//Disable slash forwarding
		controller.RedirectTrailingSlash = false

		//Add Routes
		controller.Handle("GET", "/", func(c *gin.Context) {
			c.String(http.StatusOK, "Hello sir ! (Alfred)")
		})

		// Create mocks routes
		AddMocksRoutes(controller, mocks)
	}

	return controller, nil
}

// Serve will bind the port(s) and launch serve in a separated goroutine
func Serve(main_ctx context.Context, conf *conf.Config, server *http.Server) {

	//Bind
	listener, err := net.Listen("tcp", server.Addr)
	if err != nil {
		log.Error(main_ctx, "Error Server Binding", err)
		panic("Error Server Binding on " + server.Addr)
	}

	//Log that's the bind is ok
	log.Info(main_ctx, "Alfred started to serve on host "+conf.Alfred.Core.Listen.Ip+" and is listening at port "+conf.Alfred.Core.Listen.Port)

	//Start to serve
	go func() {
		err := server.Serve(listener)
		if err != nil {
			log.Error(main_ctx, "Error at server start", err)
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

//requestRecover retreives a panic if it happened during request treatment
func requestRecover(ctx *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			ctx.AbortWithStatus(http.StatusInternalServerError)
			log.Error(ctx, "Catched Panic", errors.New(fmt.Sprint(r)))
		}
	}()

	// Process request
	ctx.Next()
}
