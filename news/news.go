/*
	Package news defines the interface for the top news API which will be used by the
	request and server modules. To use the designed server to work with the desired API
	Implement the API interface as a wrapper around the desired API.
*/
package news

import (
	"time"
)

type ID interface{}
type News interface {
	ID() ID
}

// API is an interface to a Top News API
// TODO: API must be Thread-Safe
type API interface {
	// Count() returns the number of top news required
	Count() int
	// Setup() function is called to initialize the API
	// If Setup() fails, then Server panics. Implement Initialization task here.
	Setup() error
	// Timeout() returns the maximum time the server will wait on TopNews() and News() methods
	Timeout() time.Duration
	// NumRetries() returns the number of times to retry before returning error.
	NumRetries() int
	// CacheSize() defines the capacity of LRU used for caching response.
	CacheSize() int
	// TopNews() must return the IDs of possible top news in order
	TopNews() ([]ID, error)
	// News() must return the News for a given ID
	News(ID) (News, error)
	// IsRequired() is used to check if fetched news must be in top news.
	IsRequired(News) bool

	// These are only for cached server implementation
	// which will be using this api
	// RefreshTime is the time after which cached top news will be
	// refreshed, that is, re-obtained and re-cached
	RefreshTime() time.Duration
	// FailureRefreshTime is the after which top news will be refreshed
	// if previous attempt to obtain top news failed.
	FailureRefreshTime() time.Duration
	// ExpireTime is the time after which the cached top news  will be
	// invalidated and deleted
	ExpireTime() time.Duration
	// For Good Response time, ensure that
	// refresh time + update time < expire time
	// in which case cached response will always be used
}
