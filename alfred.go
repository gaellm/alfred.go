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

package main

import (
	"alfred/internal/conf"
	"alfred/internal/log"
	"alfred/internal/mock"
	"alfred/internal/server"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

// Afred is a mock, written in Go (Golang), for performance testing. Alfred
// manages a mock list, offers helpers, permits to trigger asynchronous
// actions, and offers the ability to wrap users' javascript functions;
// users have infinite creatives possibilities.
func main() {

	configuration, err := conf.GetConfiguration()
	if err != nil {

		panic(fmt.Errorf("fatal error, config file: %w", err))
	}
	configurationJson, _ := json.Marshal(configuration)

	//Log Management
	log.InitLogger(configuration.Alfred.Name, true, configuration.Alfred.Version)

	//Core context
	ctx := context.Background()
	log.Debug(ctx, "alfred configuration initialized with: "+string(configurationJson))

	/*
		matches, err := files.FindFiles(configuration.Alfred.Core.MocksDir, "*.json")
		if err != nil {
			log.Error(ctx, "Application Panic", errors.New("error during mocks load..."+err.Error()))
			panic("error during mocks load......" + err.Error())
		}

		log.Info(ctx, "mock files loaded: "+strings.Join(matches, ","))

		mock.TestMock()
	*/

	mockCollection := mock.CreateMockCollectionFromFolder(configuration.Alfred.Core.MocksDir)
	log.Info(ctx, "mock files loaded: "+mockCollection[0].GetName())

	//------------------
	// Server Management
	//------------------
	var asyncRunningJobsCount sync.WaitGroup //Use to count process to wait before shutdowning
	classicServer, err := server.BuildServer(&configuration, &asyncRunningJobsCount)
	if err != nil {
		log.Error(ctx, "Server Panic", errors.New("error during preparing controller..."+err.Error()))
		panic("Error during preparing controller..." + err.Error())
	}

	// Let's go !!!
	server.Serve(ctx, &configuration, classicServer)

	//---------------------
	// Wait for a kill
	//---------------------
	// Set up channel on which to send signal notifications.
	// We must use a buffered channel or risk missing the signal
	// if we're not ready to receive when the signal is sent.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	// Block until a signal is received.
	<-c

	//---------------------
	// Shutdown Management
	//---------------------
	// Stop externalApiServer
	server.Stop(ctx, classicServer, &asyncRunningJobsCount)
}
