# Proxy News
This is generalisation of gophercises exercise #13 [Quiet HN](https://github.com/gophercises/quiet_hn)

## Design
Designed a generalized proxy server for Quiet HN like top news sites.
### Modules
#### news

This Package defines an API interface for Top News Websites which the designed application server uses. To used the designed server as such for other top news sites, implement the interface. The expected functionality of the methods are described below

1. Count() 
returns the number of top news required.
2. Setup()
is used to initialize the API and will be called.
3. Timeout()
returns the maximum time the server will wait for responses from TopNews() and News() methods
4. NumRetries()

