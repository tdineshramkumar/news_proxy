package request

import (
	"fmt"

	"github.com/t-drk/news_proxy/cache"
	"github.com/t-drk/news_proxy/news"
)

// SerialRequest is a function that obtained the top news using the given API. It uses the requiredCache and newsCache to
// avoid unncessary additional requests. It fetches the top news sequentially.
func SerialRequest(api news.API, requiredCache cache.Cache, newsCache cache.Cache) ([]news.News, error) {
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
	// the news list will contain the news obtained
	newsList := make([]news.News, 0, numStories)

	for _, newsID := range newsIDs {
		if newsCache.Contains(newsID) {
			if req, ok := requiredCache.Get(newsID); ok {
				// If found in required cache
				if req.(bool) {
					// If required then
					newsCached, _ := newsCache.Get(newsID)
					newsList = append(newsList, newsCached.(news.News))
				}
				// If not required then don't fetch it.
				// expiry by lru policy

			} else {
				// If not found in required cache
				// TODO: cache delete operation not implemented
				// However, the news article is fetched only if required
				// this case may result if a fetched news article become filtered later
				newsCached, _ := newsCache.Get(newsID)
				requiredStatus := api.IsRequired(newsCached.(news.News))
				requiredCache.Add(newsID, requiredStatus)
				if requiredStatus {
					newsList = append(newsList, newsCached.(news.News))
				}
				// if not required then channel will be close automatically
			}
		} else {
			// If not found in news cache
			if req, ok := requiredCache.Get(newsID); ok && !req.(bool) {
				// If not required then don't fetch
				continue
			}
			// Get the news
			newsObtained, err := ExecuteTask(
				func() (interface{}, error) {
					return api.News(newsID)
				}, timeout, numRetries)
			if err != nil {
				fmt.Println("ERROR while getting news id", newsID, "error: [", err, "]")
				// as news not obtained for the id
				continue
			}
			if api.IsRequired(newsObtained.(news.News)) {
				// If news is required, then cache the news
				newsCache.Add(newsID, newsObtained.(news.News))
				requiredCache.Add(newsID, true)
				newsList = append(newsList, newsObtained.(news.News))
			} else {
				// If news not required
				requiredCache.Add(newsID, false)
			}
		}
		// if required number of top news obtained then break
		if len(newsList) >= numStories {
			break
		}
	}
	/*
		var News news.News
		for _, id := range newsIDs {
			if requiredCache.Contains(id) {
				// If data is present in required cache
				if required, _ := requiredCache.Get(id); required.(bool) {
					// If required
					if newsCache.Contains(id) {
						// If present in news cache
						Value, _ := newsCache.Get(id)
						News = Value.(news.News)
						// Then Add it to the list
						goto ADD
					} else {
						// Then fetch it
						goto FETCH
					}
				} else {
					// If not required
					continue
				}
			}
			if newsCache.Contains(id) {
				// If not in required cache but in newsCache
				Value, _ := newsCache.Get(id)
				News = Value.(news.News)
				if !api.IsRequired(News) {
					// Add to cache as not required
					requiredCache.Add(id, false)
					continue
				}
				goto ADD
			}
		FETCH:
			value, err = ExecuteTask(func() (interface{}, error) { return api.News(id) }, timeout, numRetries)
			if err != nil {
				fmt.Println("ERROR While getting news for id ", id, err)
				continue
			}
			News = value.(news.News)
			if !api.IsRequired(News) {
				requiredCache.Add(id, false)
				continue
			}
		ADD:
			requiredCache.Add(id, true)
			newsCache.Add(id, News)
			newsList = append(newsList, News)
			if len(newsList) >= numStories {
				break
			}
		}
	*/
	if len(newsList) == 0 {
		return nil, NoNewsError("No News Obtained.")
	}
	return newsList, nil
}
