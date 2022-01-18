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
	log.InitLogger("alfred", true, "0.1")

	//Core context
	ctx := context.Background()
	log.Debug(ctx, "alfred configuration initialized with: "+string(configurationJson))

	matches, err := files.FindFiles(configuration.Core.MocksDir, "*.json")
	if err != nil {
		log.Error(ctx, "Application Panic", errors.New("error during mocks load..."+err.Error()))
		panic("error during mocks load......" + err.Error())
	}

	log.Info(ctx, "mock files loaded: "+strings.Join(matches, ","))
}
