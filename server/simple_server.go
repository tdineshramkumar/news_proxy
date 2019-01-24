package server

import (
	"fmt"

	"github.com/t-drk/news_proxy/cache"
	"github.com/t-drk/news_proxy/lru"
	"github.com/t-drk/news_proxy/news"
	"github.com/t-drk/news_proxy/request"
)

type simpleServer struct {
	parallelize   bool
	numRetries    int
	requiredCache cache.Cache
	newsCache     cache.Cache
	api           news.API
}

// SimpleServer returns a server implementation that does not cache server responses.
// api is the implemention of the news API
// parallelize indicates whether serial or parallel requests.
// numRetries is used only in case of parallel requests.
// Note: news items are lru cached to avoid unnecessary fetching.
func SimpleServer(api news.API, parallelize bool, numRetries int) Server {
	ss := new(simpleServer)
	ss.parallelize = parallelize
	ss.numRetries = numRetries
	ss.requiredCache = lru.LRU_TS(api.CacheSize())
	ss.newsCache = lru.LRU_TS(api.CacheSize())
	ss.api = api
	return ss
}

func (ss *simpleServer) HandleRequest() []news.News {
	var (
		topNews []news.News
		err     error
	)
	if ss.parallelize {
		topNews, err = request.ParallelRequest(ss.api, ss.requiredCache, ss.newsCache, ss.numRetries)
	} else {
		topNews, err = request.SerialRequest(ss.api, ss.requiredCache, ss.newsCache)
	}
	if err != nil {
		// if error while obtaining top news.
		fmt.Println("ERROR simple server while obtaining top news. error :[", err, "]")
		return nil
	}
	return topNews

}
