package lru

import "github.com/t-drk/news_proxy/cache"

/*
	Makes the lru module thread safe by serializing the operations using a channel.
*/

type addStruct struct {
	key   cache.KEY
	value cache.VALUE
}
type getStruct struct {
	key      cache.KEY
	response chan getRespStruct
}
type containsStruct struct {
	key      cache.KEY
	response chan containsRespStruct
}
type getRespStruct struct {
	value cache.VALUE
	ok    bool
}
type containsRespStruct struct {
	ok bool
}

type lru_ts struct {
	cache    cache.Cache
	add      chan addStruct
	get      chan getStruct
	contains chan containsStruct
}

// TSServer is a server goroutine that handles requests to the underlying lru cache.
// It processes the requests serially to avoid race condition.
// TODO: Make the server similar to a reader-writer model, as reads can happen in parallel.
func TSServer(Cache *lru_ts) {
	for {
		select {
		case add := <-Cache.add:
			// IF ADD request to cache
			Cache.cache.Add(add.key, add.value)
		case get := <-Cache.get:
			// IF GET request to cache
			value, ok := Cache.cache.Get(get.key)
			get.response <- getRespStruct{value, ok}
		case contains := <-Cache.contains:
			// IF CONTAIN request to cache
			ok := Cache.cache.Contains(contains.key)
			contains.response <- containsRespStruct{ok}
		}
	}
}

// LRU Creates and returns LRU that is thread safe
func LRU_TS(capacity int) cache.Cache {
	Cache := new(lru_ts)
	Cache.cache = LRU(capacity)
	Cache.add = make(chan addStruct)
	Cache.get = make(chan getStruct)
	Cache.contains = make(chan containsStruct)
	go TSServer(Cache)
	return Cache
}

/*
	The operations are send to a channel for processing by the server.
	It blocks till operations complete to ensure consistency.
*/
func (c *lru_ts) Add(key cache.KEY, value cache.VALUE) {
	c.add <- addStruct{key, value}
}

func (c *lru_ts) Get(key cache.KEY) (value cache.VALUE, ok bool) {
	response := make(chan getRespStruct, 1)
	c.get <- getStruct{key, response}
	getResp := <-response
	return getResp.value, getResp.ok
}

func (c *lru_ts) Contains(key cache.KEY) bool {
	response := make(chan containsRespStruct, 1)
	c.contains <- containsStruct{key, response}
	containsResp := <-response
	return containsResp.ok
}
