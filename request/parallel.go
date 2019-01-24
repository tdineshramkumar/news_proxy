package request

import (
	"fmt"

	"github.com/t-drk/news_proxy/cache"
	"github.com/t-drk/news_proxy/news"
)

type NoNewsError string

func (err NoNewsError) Error() string {
	return string(err)
}

// ParallelRequest is a function that obtains the top news using the given API. It uses the requiredCache and newsCache to
// avoid unnecessary additional requests. This function requests the top news concurrently as dictated by the numRoutines parameter..
func ParallelRequest(api news.API, requiredCache cache.Cache, newsCache cache.Cache, numRoutines int) ([]news.News, error) {
	if numRoutines <= 0 {
		panic(fmt.Sprintf("ParallelRequest function got unexceptable number of routines.[%d]", numRoutines))
	}

	timeout := api.Timeout()
	numRetries := api.NumRetries()
	numStories := api.Count()
	newsIDsObtained, err := ExecuteTask(func() (interface{}, error) { return api.TopNews() }, timeout, numRetries)
	if err != nil {
		fmt.Println("ERROR while getting top news ids. error:[", err, "]")
		return nil, err
	}
	// the news IDs obtained
	newsIDs := newsIDsObtained.([]news.ID)
	// newsIDChan is a channel that will send requests to goroutines which will fetch
	// the news concurrently. it follows a request response structure.
	// and it is buffered channel.
	newsIDChan := make(chan struct {
		// the requested news ID
		ID news.ID
		// the channel through which response is send
		newsChan chan news.News
	}, len(newsIDs))

	// newsChannels is a list of response channels which will be sorted based on the
	// order of top news ids.
	// all the channels currently are nil channels.
	newsChannels := make([]chan news.News, len(newsIDs))
	// for each of the top news id create a buffered response channel for non-blocking
	// response.
	for i, id := range newsIDs {
		newsChan := make(chan news.News, 1)
		newsIDChan <- struct {
			ID       news.ID
			newsChan chan news.News
		}{id, newsChan}
		// obtain response in order
		newsChannels[i] = newsChan
	}
	// close the news ID channel
	close(newsIDChan)
	// the quit channel which indicates the completion of operation
	// indicates that no more news needs to be fetched
	quitChannel := make(chan bool)
	// newsList := make([]news.News, 0, numStories)

	// the function to fetch id and send it to the response channel
	// or close the channel if not necessary
	fetchID := func(newsID news.ID, resp chan news.News) {
		// fmt.Println("Fetching news ID", newsID)
		defer close(resp)
		if newsCache.Contains(newsID) {
			if req, ok := requiredCache.Get(newsID); ok {
				// If found in required cache
				if req.(bool) {
					// If required then
					newsCached, _ := newsCache.Get(newsID)
					resp <- newsCached.(news.News)
				}
				// If not required then don't fetch it.
				// expiry by lru policy
				return

			} else {
				// If not found in required cache
				// TODO: cache delete operation not implemented
				// However, the news article is fetched only if required
				// this case may result if a fetched news article become filtered later
				newsCached, _ := newsCache.Get(newsID)
				requiredStatus := api.IsRequired(newsCached.(news.News))
				requiredCache.Add(newsID, requiredStatus)
				if requiredStatus {
					resp <- newsCached.(news.News)
				}
				// if not required then channel will be close automatically
				return
			}
		} else {
			// If not found in news cache
			if req, ok := requiredCache.Get(newsID); ok && !req.(bool) {
				// If not required then don't fetch
				return
			}
			// Get the news
			newsObtained, err := ExecuteTask(
				func() (interface{}, error) {
					return api.News(newsID)
				}, timeout, numRetries)
			if err != nil {
				fmt.Println("ERROR while getting news id", newsID, "error: [", err, "]")
				// as news not obtained for the id
				return
			}
			if api.IsRequired(newsObtained.(news.News)) {
				// If news is required, then cache the news
				newsCache.Add(newsID, newsObtained.(news.News))
				requiredCache.Add(newsID, true)
				resp <- newsObtained.(news.News)
			} else {
				// If news not required
				requiredCache.Add(newsID, false)
			}
		}
	}

	// lauch go routines to obtain the news concurrently
	for routine := 0; routine < numRoutines; routine++ {
		go func() {
			for {
				select {
				// IF top news obtained the quitChannel becomes readable
				case <-quitChannel:
					return
					// Obtain the request
				case r, ok := <-newsIDChan:
					if ok {
						fetchID(r.ID, r.newsChan)
						continue
					}
					return
				}
			}
		}()
	}
	/*
		fmt.Println("Waiting for responses")
		fmt.Println("Response Channels")
		for i, newsChan := range newsChannels {
			fmt.Println("i:", i, "newsChan:", newsChan, "len:", len(newsChan), cap(newsChan))
		}
	*/
	// Make a news slice which at present can contain the num of stories requested
	newsList := make([]news.News, 0, numStories)
	// read the responses
	for _, newsChan := range newsChannels {
		newsObtained, ok := <-newsChan
		// fmt.Println("News Channel", i, "returned", ok, "total", len(newsChannels))
		if ok {
			// If news not filtered
			newsList = append(newsList, newsObtained)
		}
		if len(newsList) >= numStories {
			// If requested number of news channels obtained
			break
		}
	}
	close(quitChannel)
	if len(newsList) == 0 {
		// if no news obtained
		return nil, NoNewsError("No News Obtained.")
	}
	// return the news obtained.
	return newsList, nil
}
