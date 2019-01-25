package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
	"time"
)

func main() {
	var port, numTrials, numRoutines int
	var serverDomain string
	flag.IntVar(&port, "port", 9999, "port number of webserver running in localhost.")
	flag.IntVar(&numTrials, "trials", 1, "number of requests to send per routine.")
	flag.IntVar(&numRoutines, "routines", 1, "number of concurrent routines to send request.")
	flag.StringVar(&serverDomain, "server", "localhost", "domain name of the server to test.")
	flag.Parse()

	URL := fmt.Sprintf("http://%s:%d/", serverDomain, port)
	var wg sync.WaitGroup

	// Customize the Transport to have larger connection pool
	defaultRoundTripper := http.DefaultTransport
	defaultTransportPointer, ok := defaultRoundTripper.(*http.Transport)
	if !ok {
		panic(fmt.Sprintf("defaultRoundTripper not an *http.Transport"))
	}
	defaultTransport := *defaultTransportPointer // dereference it to get a copy of the struct that the pointer points to
	defaultTransport.MaxIdleConns = 100
	if numRoutines > 100 {
		defaultTransport.MaxIdleConnsPerHost = numRoutines
	} else {
		defaultTransport.MaxIdleConnsPerHost = 100
	}
	myClient := &http.Client{Transport: &defaultTransport}

	for routine := 0; routine < numRoutines; routine++ {
		wg.Add(1)
		go func() {

			for trial := 0; trial < numTrials; trial++ {
				start := time.Now()
				// Send the request to the server
				resp, err := myClient.Get(URL)
				io.Copy(ioutil.Discard, resp.Body)
				resp.Body.Close()
				duration := time.Now().Sub(start)
				responseTime := fmt.Sprintf("%f", duration.Seconds())
				if err != nil {
					fmt.Fprintln(os.Stderr, "[", responseTime, "]", err)
					continue
				}
				fmt.Fprintln(os.Stdout, responseTime)
			}
			wg.Done()
		}()
	}
	wg.Wait()
}
