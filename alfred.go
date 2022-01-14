package main

import (
	"alfred/internal/conf"
	"alfred/pkg/files"
	"fmt"
)

// Afred is a mock, written in Go (Golang), for performance testing. Alfred
// manages a mock list, offers helpers, permits to trigger asynchronous
// actions, and offers the ability to wrap users' javascript functions;
// users have infinite creatives possibilities.
func main() {

	// Init Alfred configuration.
	config, err := conf.InitConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	matches, err := files.FindFiles(config.GetString(conf.MOCKS_DIR), "*.json")
	if err != nil {
		panic(fmt.Errorf("no mock files found: %w", err))
	}

	fmt.Println(matches)

}
