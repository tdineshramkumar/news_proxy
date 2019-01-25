/*
	Hackker News API is an Implementation of NEW API.
*/
package hn

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/t-drk/news_proxy/news"
)

const (
	apiBase = "https://hacker-news.firebaseio.com/v0"
)

// Item represents a single item returned by the HN API. This can have a type
// of "story", "comment", or "job" (and probably more values), and one of the
// URL or Text fields will be set, but not both.
//
// FOr the purpose of this exercise, we only care about items where the
// type is "story:, and the URL is set.
type Item struct {
	By          string `json:"by"`
	Descendants int    `json:"descendants"`
	Id          int    `json:"id"`
	Kids        []int  `json:"kids"`
	Score       int    `json:"score"`
	Time        int    `json:"time"`
	Title       string `json:"title"`
	Type        string `json:"type"`

	// Only one of these should exist
	Text string `json:"text"`
	URL  string `json:"url"`

	Host string
}

func (item *Item) ID() news.ID {
	return item.Id
}

type API struct {
	topNews    int
	timeout    time.Duration
	numRetries int
	cacheSize  int
	client     *http.Client
}

func New() *API {
	timeout := time.Second * 20
	api := API{topNews: 30, timeout: timeout, numRetries: 3, cacheSize: 1000}
	return &api
}
func (api *API) Count() int {
	return api.topNews
}

func (api *API) Setup() error {
	// Customize the Transport to have larger connection pool
	defaultRoundTripper := http.DefaultTransport
	defaultTransportPointer, ok := defaultRoundTripper.(*http.Transport)
	if !ok {
		panic(fmt.Sprintf("defaultRoundTripper not an *http.Transport"))
	}
	defaultTransport := *defaultTransportPointer // dereference it to get a copy of the struct that the pointer points to
	defaultTransport.MaxIdleConns = 100
	defaultTransport.MaxIdleConnsPerHost = 100

	api.client = &http.Client{Transport: &defaultTransport}
	return nil
}

func (api *API) Timeout() time.Duration {
	return api.timeout
}

func (api *API) NumRetries() int {
	return api.numRetries
}

func (api *API) CacheSize() int {
	return api.cacheSize
}

func (api *API) TopNews() (ids []news.ID, err error) {
	resp, err := api.client.Get(fmt.Sprintf("%s/topstories.json", apiBase))
	if err != nil {
		fmt.Println("ERROR Could not get HN TopNews", err)
		return nil, err
	}
	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&ids)
	if err != nil {
		fmt.Println("ERROR Could not decode HN TopNews IDs", err)
		return nil, err
	}
	return
}

func (api *API) News(ID news.ID) (news.News, error) {
	var item Item
	id := int(ID.(float64))
	resp, err := api.client.Get(fmt.Sprintf("%s/item/%d.json", apiBase, id))
	if err != nil {
		fmt.Println("ERROR while getting news", err)
		return nil, err
	}
	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&item)
	if err != nil {
		fmt.Println("ERROR while decoding item", err)
		return nil, err
	}
	url, err := url.Parse(item.URL)
	if err == nil {
		item.Host = strings.TrimPrefix(url.Hostname(), "www.")
	} else {
		fmt.Println("ERROR While parsing URL", err)
	}
	return &item, err
}

func (api *API) IsRequired(news news.News) bool {
	var item *Item = news.(*Item)
	return item.Type == "story" && item.URL != ""
}

func (api *API) RefreshTime() time.Duration {
	//	return time.Second * 10
	return time.Minute * 10
}

func (api *API) ExpireTime() time.Duration {
	// return time.Second * 20
	return time.Minute * 12
}

func (api *API) FailureRefreshTime() time.Duration {
	// return time.Second * 5
	return time.Minute * 2
}
