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

### Designs
#### Top News API
First Setup method is called to initialize the API. 
For each request for top news by the client, first top news IDs are obtained by TopNews() method, then news for each of the IDs is obtained using the News(ID) method. The number of news requested depend on the number of top news requested and number of news items that are filtered, it also depends on whether requests are sent serially or parallely. Parallel requests tend to over fetch.

#### Parallel Request 
Once top news IDs are obtained via a request to hn, it creates a buffered channel for requests for news IDs, a request is a structure containing of the news ID and a response buffered channel, via which news is obtained back in order. It then creates a pool of goroutines, which read requests from the request channel, it obtains the news and checks if matches the filter, if yes then sends it via the response channel, which was part of the request, else it closes the response channel for that request to indicate that news was filtered. It then fetches the next request from the request channel. This process continues till either there are no requests or a quit message is sent to indicate top news were obtained. (Quit message is send via another channel once required number of top news were obtained). Now responses of the goroutine pool are read in order, if response channel is closed, implying filtered news, it reads next response, it keeps reading responses till required number of top news is obtained, then it sends a quit message to the thread pool by closing the quit channel.

#### Cached Server
A goroutine waits for requests from client, which contain a buffered response channel on which cached server responds. Cached server is designed to increase throughput by processing as much events as possible which would not block or stall the routine. It asynchronously and periodically sends request for top news and refreshes the cached top news dictated by the refresh time and expiry time. It only listen for client requests only when cached responses are available. It fetches top news periodically even when no clients sends any request to maintain the freshness of the cache.

