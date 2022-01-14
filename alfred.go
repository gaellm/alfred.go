package main

import (
	"alfred/collection"
	"fmt"
)

// Afred is a mock, written in Go (Golang), for performance testing. Alfred manages a mock list, offers helpers,
// permits to trigger asynchronous actions, and offers the ability to wrap users' javascript functions; users have
// infinite creatives possibilities.
func main() {

	matches, err := collection.FindFiles("resources/user-files/mocks/", "*.jssdon")
	if err != nil {
		panic("no mock files found")
	}

	fmt.Println(matches)
}
