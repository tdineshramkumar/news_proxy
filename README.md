# Proxy News
This is generalisation of gophercises exercise #13 [Quiet HN](https://github.com/gophercises/quiet_hn)

### Modules
#### news

This Package defines an API interface for Top News Websites which the designed application server uses. To used the designed server as such for other top news sites, implement the interface. 

#### request 

This Package defines ExecuteTask, SerialRequest and ParallelRequest functions. ExecuteTask function executes a given function with a given timeout. If execution time exceeds timeout, then it retries for a maximum of specified retries. SerialRequest and ParallelRequest fetches the top news serially and parallely respectively. Both of them first fetch the top news IDs. They defer in how they fetch the news items for those IDs. SerialRequest fetches them one after another, checks if they match the required filter and fetches till required number of top news is obtained. ParallelRequest fetches the news items using a thread pool (or rather a goroutine pool) of a specified size. It then arranges them back in order to obtain specified number of top news in order.

#### cache
cache package defines an interface for cache with Add, Contains and Get method.

#### lru
lru package implements the cache interface with least recently used policy for eviction when entries exceed the capacity. It also contains a thread safe implementation of the lru, which achieves the same by serializing the requests to a single server which maintains the lru.


#### server
server package defines interface for a server which returns top news. It defines a single method which will return the top news. server package have two implementations, one which requests calls the top news API for each request and another which caches the top news periodically and returns the cached top news to improve response time and increase throughput of the server.

#### hn
hn package implements the news API interface for Quick HN.

