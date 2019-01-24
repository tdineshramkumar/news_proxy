package main

import (
	"flag"
	"fmt"
	"github.com/t-drk/news_proxy/hn"
	"github.com/t-drk/news_proxy/news"
	"github.com/t-drk/news_proxy/server"
	"log"
	"net/http"
	"text/template"
	"time"
)

func main() {
	var parallelize, enableCaching bool
	flag.BoolVar(&parallelize, "parallel", false, "parallelize the requests to HN server")
	flag.BoolVar(&enableCaching, "caching", false, "create a caching application server")
	flag.Parse()
	// Obtain the template.
	tpl := template.Must(template.ParseFiles("index.gohtml"))
	// Obtain the HackerNews API.
	API := hn.New()
	// Create a server using the required flags
	Server := server.New(API, parallelize, enableCaching)
	/*
		var Server server.Server
		switch {
		case parallelize && enableCaching:
			fmt.Println("Creating server with parallelize and caching.")
			Server = server.New(API, )
		case parallelize && !enableCaching:
			fmt.Println("Creating server with parallelize.")
		case !parallelize && enableCaching:
			fmt.Println("Creating server with caching.")
		case !parallelize && !enableCaching:
			fmt.Println("Creating server with no parallelize and no caching.")
		}
	*/
	// Create a http server which will interact with this above application server which interacts with the APIs
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		topNews := Server.HandleRequest()
		if topNews == nil {
			// If Server Error
			http.Error(w, "Failed to load top stories", http.StatusInternalServerError)
			return
		}
		err := tpl.Execute(w, templateData{
			topNews,
			time.Now().Sub(startTime),
		})
		if err != nil {
			log.Fatal(fmt.Sprintln("tpl.Execute returned error: [", err, "]"))
		}
	})

	log.Fatal(http.ListenAndServe(":9999", nil))
	/*
		tpl := template.Must(template.ParseFiles("index.gohtml"))
		api := hn.New()
		if err := api.Setup(); err != nil {
			// If failed to setup API
			panic(err.Error())
		}
		requiredCache, newsCache := lru.LRU_TS(api.CacheSize()), lru.LRU_TS(api.CacheSize())
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			startTime := time.Now()
			newsList, err := request.SerialRequest(api, requiredCache, newsCache)
			// newsList, err := request.ParallelRequest(api, requiredCache, newsCache)
			if err != nil {
				fmt.Println("ERROR whil serial request", err)
				http.Error(w, "Failed to load top stories", http.StatusInternalServerError)
			} else {
				responseTime := time.Now().Sub(startTime)
				stories := newsList // .([]hn.Item)
				data := templateData{stories, responseTime}
				err = tpl.Execute(w, data)
				if err != nil {
					fmt.Println(err)
					http.Error(w, "Failed to process template", http.StatusInternalServerError)
				}
			}
		})

		// start the server
		log.Fatal(http.ListenAndServe(":9999", nil))
	*/
}

type templateData struct {
	Stories []news.News
	Time    time.Duration
}
