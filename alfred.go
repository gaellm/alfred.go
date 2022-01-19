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
	"alfred/pkg/files"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
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

	matches, err := files.FindFiles(configuration.Alfred.Core.MocksDir, "*.json")
	if err != nil {
		log.Error(ctx, "Application Panic", errors.New("error during mocks load..."+err.Error()))
		panic("error during mocks load......" + err.Error())
	}

	log.Info(ctx, "mock files loaded: "+strings.Join(matches, ","))
}
