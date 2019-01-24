/*
Package cache defines an interface for cache. Cache is key, value pair store.
key and value used in the cache have the type of empty interface. Use type assertions
while retrieving from the cache.
Add method adds a key value pair to cache
Contains method checks if a key exists in the cache
Get returns the value stored in the cache for the given key if it exists
*/
package cache

// KEY and VALUE denote types of pair stored in cache
type (
	KEY   interface{}
	VALUE interface{}
)

// Cache is an interface for a cache
// Add adds a new pair
// Contains, Get check if element exists in cache
type Cache interface {
	Add(KEY, VALUE)
	Contains(KEY) bool
	Get(KEY) (VALUE, bool)
}

// TODO: Later add Delete method which will remove a key value pair from the cache.
