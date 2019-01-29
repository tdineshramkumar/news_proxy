package server

import (
	"sync"
	"time"

	"github.com/t-drk/news_proxy/lru"
	"github.com/t-drk/news_proxy/news"
	"github.com/t-drk/news_proxy/pool"
	"github.com/t-drk/news_proxy/request"
)

type HandleRequest func() []news.News

func (f HandleRequest) HandleRequest() []news.News {
	return f()
}

const (
	SERVER_POOL_SIZE    = 100
	REQUEST_BUFFER_SIZE = 1000
)

func loadData() {

}
func NewCachedServerv2(api news.API) Server {
	requiredCache := lru.LRU_TS(api.CacheSize())
	newsCache := lru.LRU_TS(api.CacheSize())
	threadPool := pool.New(api.PoolSize())

	requestsChannel := make(chan chan []news.News, REQUEST_BUFFER_SIZE)
	var writeLock sync.RWMutex
	cachedResponse, _ := request.SerialRequest(api, requiredCache, newsCache, threadPool)

	// This goroutine periodically updates the cache
	go func() {
		for range time.Tick(api.RefreshTime()) {
			// Get the latest news
			newResponse, _ := request.SerialRequest(api, requiredCache, newsCache, threadPool)
			func() {
				// Update the cache
				writeLock.Lock()
				defer writeLock.Unlock()
				cachedResponse = newResponse
			}()
		}
	}()

	for i := 0; i < SERVER_POOL_SIZE; i++ {
		go func() {
			for responseChannel := range requestsChannel {
				func() {
					writeLock.RLock()
					defer writeLock.RUnlock()
					responseChannel <- cachedResponse
				}()
			}
		}()
	}

	return HandleRequest(func() []news.News {
		responseChannel := make(chan []news.News, 1)
		requestsChannel <- responseChannel
		return <-responseChannel
	})
}
