// package server defines application servers which hanwhich handle requests differently
package server

import (
	"fmt"
	"github.com/t-drk/news_proxy/news"
)

// Server is an interface for a server implementation. It defines a single method to handle request
// for client. HandleRequest method returns the top news or nil if internal server error occurs.
type Server interface {
	HandleRequest() []news.News
}

/*
// Flag is an enumeration to define the type of server that will be returned by New function.
type Flag int

// By Default no flags are set
// PARALLELIZE indicates that server should fetch top news items in parallel
// CACHE indicates that server should cache the top news (server caches the top news for subsequent queries.)
const (
	PARALLELIZE Flag = 10
	CACHE       Flag = 20
)
*/
// New function returns a Server whose type is defined by the flags.
// It takes the Top News API as input.

// func New(api news.API, flags ...Flag) Server {
func New(api news.API, parallelize, enableCaching bool) Server {
	// Initialize the API
	err := api.Setup()
	if err != nil {
		// IF API Initialization fails then panic
		panic(fmt.Sprintf("News API Setup Failed with error `%v`", err))
	}
	/*
		// By default caching responses and parallelization of requests are not enabled
		var parallelize, enableCaching = false, false
		fmt.Println("Server Created with", len(flags), "flags.")
		for flag := range flags {
			switch flag {
			case PARALLELIZE:
				fmt.Println("Parallelization is enabled.")
				parallelize = true
			case CACHE:
				fmt.Println("Caching is enabled.")
				enableCaching = true
			}
		}
	*/
	// Num Retries currently defined as constant based on number of top news.
	numRetries := int(float64(api.Count()) * 1.5)
	if enableCaching {
		// return the caching server.
		fmt.Println("Launching a Cached Server. Parallelize:", parallelize)
		return NewCachedServer(api, parallelize, numRetries)
	}
	fmt.Println("Launching a simple server. Parallelize:", parallelize)
	// Return simple server if no caching is enabled
	return SimpleServer(api, parallelize, numRetries)
}
