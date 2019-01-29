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

func loadData() {

}
func NewCachedServerv2(api news.API) Server {
	requiredCache := lru.LRU_TS(api.CacheSize())
	newsCache := lru.LRU_TS(api.CacheSize())
	threadPool := pool.New(api.PoolSize())

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

	return HandleRequest(func() []news.News {
		writeLock.RLocker()
		defer writeLock.RUnlock()
		return cachedResponse
	})
}
