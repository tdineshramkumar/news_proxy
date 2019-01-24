package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"
)

func main() {
	var port, numTrials int
	flag.IntVar(&port, "port", 9999, "webserver port number")
	flag.IntVar(&numTrials, "trials", 1, "number of requests to send")
	flag.Parse()

	URI := fmt.Sprintf("http://localhost:%d/", port)
	for i := 0; i < numTrials; i++ {
		start := time.Now()
		// Send the request to the server
		_, err := http.Get(URI)
		responseTime := time.Now().Sub(start)
		outputString := fmt.Sprintf("%f", responseTime.Seconds())
		if err != nil {
			// some error occurred.
			fmt.Fprintln(os.Stderr, outputString)
			continue
		}
		fmt.Fprintln(os.Stdout, outputString)
	}
}
