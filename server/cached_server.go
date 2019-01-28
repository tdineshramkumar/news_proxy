package server

import (
	"fmt"
	"time"

	"github.com/t-drk/news_proxy/lru"
	"github.com/t-drk/news_proxy/news"
	"github.com/t-drk/news_proxy/pool"
)

// cachedServer is like the simpleServer however caches the response it obtains on one request to
// send it to another request. THus improving response time. However it results in increased complexity.
// requestChannel is the interface to the caching server which maintains the cached responses.
type cachedServer struct {
	simpleServer
	requestChannel chan chan []news.News
}

// getNews wraps around the HandleRequest method of the simple server. calls it asynchronously and sends the
// response through a buffered channel (for non-blocking sender design)
// we try to make the caching server as fast as possible, reducing blocks whenever possible.
func (cs *cachedServer) getNews() chan []news.News {
	newsChannel := make(chan []news.News, 1)
	go func() { newsChannel <- cs.simpleServer.HandleRequest() }()
	return newsChannel
}

// cachingServer is the main part of the cached server, which maintains the cached responses. It receives requests
// from a channel which returns another channel through which it responds. Caching Server does not respond to requests
// when its cache is invalidated. It is designed to block minimally and such that response time is improved (by caching)
// Also it is designed to maximize throughput by asynchronous behaviour.
func (cs *cachedServer) cachingServer() {
	// cachedNews will contained the cached top news
	// refreshTimeout, expiredTimeout are the timeout channels
	// newsChannel is the channel for asynchronous news responses.
	var (
		cachedNews                    []news.News
		refreshTimeout, expireTimeout <-chan time.Time
		newsChannel                   = cs.getNews()
		// newsValid indicates if the news is valid or not
		newsValid = false
	)
	for {
		requestChannel := cs.requestChannel
		if !newsValid {
			// In case top news not cached then, dont listen for input requests
			requestChannel = nil
		}
		select {
		case cachedNews = <-newsChannel:
			// on reading invalidate the news channel
			newsChannel = nil
			// validate the news
			newsValid = true
			if cachedNews == nil {
				// Note in case of failure to obtain data use the failure refresh time
				refreshTimeout = time.After(cs.api.FailureRefreshTime())
				expireTimeout = time.After(cs.api.ExpireTime())
				fmt.Println("Attempt to update cache failed. Will try again after FailureRefreshTime")
				continue
			}
			refreshTimeout, expireTimeout = time.After(cs.api.RefreshTime()), time.After(cs.api.ExpireTime())
			fmt.Println("Updated cached news.")
		case responseChannel := <-requestChannel:
			// in case client requests for the cached news and
			// news is cachd then send the cached news.
			// since we design for fast responses, it is expected that
			// responseChannel is a buffered channel
			responseChannel <- cachedNews
		case <-refreshTimeout:
			// Refresh the cache without invalidating the existing cache
			refreshTimeout = nil
			// However if any existing request in progress do not resend
			// If however request takes too long, then it will refreshed once cache  is invalidated on expire timeout
			if newsChannel != nil {
				// If any request in process
				fmt.Println("Cache refresh cancelled as existing request in progress.")
				// Also create a refresh timeout, this only serves the purpose of logging the above message in case if the news channel didn;t get free yet.
				// It will be invalidated if request succeeds.
				refreshTimeout = time.After(cs.api.RefreshTime())
				continue
			}
			newsChannel = cs.getNews()
			fmt.Println("Sending request to refresh cached news.")
		case <-expireTimeout:
			// InValidate the cache and send a request to update the cache.
			newsValid = false
			expireTimeout = nil
			refreshTimeout = nil
			newsChannel = cs.getNews()
			fmt.Println("Cached News expired and flushed out. Waiting for latest top news.")

		}
	}

}

func NewCachedServer(api news.API, parallelize bool, numRetries int) Server {
	cs := new(cachedServer)
	cs.parallelize = parallelize
	cs.numRetries = numRetries
	// create thread safe caches.
	cs.requiredCache = lru.LRU_TS(api.CacheSize())
	cs.newsCache = lru.LRU_TS(api.CacheSize())
	cs.api = api
	cs.requestChannel = make(chan chan []news.News)
	cs.threadPool = pool.New(api.PoolSize())
	// launch the caching server which will take care of caching the top news.
	go cs.cachingServer()
	return cs
}

func (cs *cachedServer) HandleRequest() []news.News {
	// Make a channel on which to receive the response
	// And make it a buffered one.
	responseChannel := make(chan []news.News, 1)
	cs.requestChannel <- responseChannel
	// wait on the response channel for the response from the caching server.
	topNews := <-responseChannel
	fmt.Println("Obtained top news from response caching server.")

	return topNews
}
