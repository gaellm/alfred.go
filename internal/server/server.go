package server

import (
	"alfred/internal/conf"
	"alfred/internal/log"
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

//Build the service
func BuildServer(conf *conf.Config, asyncRunningJobsCount *sync.WaitGroup) (*http.Server, error) {
	//Build all endpoints handler
	handler, err := BuildHandler(conf, asyncRunningJobsCount)
	if err != nil {
		return nil, err
	}

	//Associate the handler to a server (-> contains listening interface(s))
	//Here, the server is listening on ALL interfaces and binding on 'conf.Port' port
	return &http.Server{
		Handler: handler,
		Addr:    fmt.Sprintf("0.0.0.0:%v", 8080),
	}, nil
}

//Create Handler with endpoints to serve: used from service or from tests
func BuildHandler(conf *conf.Config, asyncRunningJobsCount *sync.WaitGroup) (http.Handler, error) {
	//log
	log.Info(context.Background(), "Starting")
	defer log.LogPanic()

	// Create controller
	gin.SetMode(gin.ReleaseMode)
	controller := gin.Default()

	{
		//Disable slash forwarding
		controller.RedirectTrailingSlash = false

		//Add Routes
		controller.GET("/", func(c *gin.Context) {
			c.String(http.StatusOK, "Hello")
		})
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
	log.Info(main_ctx, "Started Server")

	//Start to serve
	go func() {
		server.Serve(listener)
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
