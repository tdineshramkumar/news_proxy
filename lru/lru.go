/*
	Package lru contains Least Recently Used implement of Cache Interface.
	LRU defines a least recently used cache
	LRU_TS makes the LRU implementation of Cache interface thread safe suitable for concurrent access/
	lru is implemented as a doubly linked list.
*/
package lru

import "github.com/t-drk/news_proxy/cache"

type node struct {
	prev, next *node
	key        cache.KEY
}

// lru is data structure for maintaining the lru cache
// capacity denotes the maximum number of entries in the cache
// store contains the actual data (key, value) pair
// nodes maintains pointers to node(s), the linked list representing the lru.
// most and least denote the most and least recently used elements
// of the cache
type lru struct {
	capacity int
	store    map[cache.KEY]cache.VALUE
	nodes    map[cache.KEY]*node
	most     *node
	least    *node
}

// NIL function returns nil cache
// This function is used to create cache when caching is not required but function demands them as arguments.
// lru  methods handles nil values.
func NIL() cache.Cache {
	return nil
}

// LRU creates and returns an Least Recently Used Cache of given capacity
func LRU(capacity int) cache.Cache {
	if capacity <= 0 {
		panic("Capacity of LRU cache is less than 0")
	}
	Cache := new(lru)
	Cache.store = make(map[cache.KEY]cache.VALUE)
	Cache.nodes = make(map[cache.KEY]*node)
	Cache.most, Cache.least = nil, nil // not required though
	Cache.capacity = capacity
	return Cache
}

// Add method adds a key, value pair to the LRU cache
func (c *lru) Add(key cache.KEY, value cache.VALUE) {
	if c == nil {
		return
	}
	// add to store the key, value pair
	c.store[key] = value
	// Now update the LRU cache structure
	// delete from store if necessary
	if n, ok := c.nodes[key]; ok {
		// If node already present in the cache
		// Nothing to delete, just update the pointers
		if n == c.most {
			// don't need to anything
		} else if n == c.least {
			// if least and not most ie next element exist
			next, most := n.next, c.most
			next.prev, most.next = nil, n
			n.next, n.prev = nil, most
			c.least, c.most = next, n
		} else {
			// some nodes exist before and after
			prev, next, most := n.prev, n.next, c.most
			next.prev, prev.next, most.next = prev, next, n
			n.next, n.prev = nil, most
			c.most = n
		}
		return
	}
	// If node does not already exist in the cache
	n := new(node)
	n.key = key
	if len(c.nodes) < c.capacity {
		// If cache is not full
		c.nodes[key] = n
		if len(c.nodes) == 1 {
			// if no nodes in cache previously
			c.most, c.least = n, n
			n.next, n.prev = nil, nil
			return
		}
		// If at least one node in cache
		most := c.most
		most.next = n
		n.next, n.prev = nil, most
		c.most = n
		return
	} else {
		c.nodes[key] = n
		// If cache is already full
		if c.capacity <= 0 {
			panic("Capacity of cache must be greater than 0")
		}
		if c.capacity == 1 {
			// if capacity is 1 then replace existing element with current one.
			most := c.most
			k := most.key
			delete(c.store, k)
			delete(c.nodes, k)
			c.most, c.least = n, n
			return
		}
		// If capacity is greater than 1
		most, least := c.most, c.least
		k := least.key
		first := least.next
		first.prev, most.next = nil, n
		n.next, n.prev = nil, most
		c.most, c.least = n, first
		delete(c.store, k)
		delete(c.nodes, k)
		return
	}

}

// Get method returns the value if key exists and also updates most recently accessed
func (c *lru) Get(key cache.KEY) (cache.VALUE, bool) {
	if c == nil {
		var value cache.VALUE
		return value, false
	}
	value, ok := c.store[key]
	if ok {
		// add back the value so as to update the LRU to make the element the least recently used.
		c.Add(key, value)
	}
	return value, ok
}

// Contains method checks if in cache however does not update the cache.
func (c *lru) Contains(key cache.KEY) bool {
	if c == nil {
		return false
	}
	_, ok := c.store[key]
	return ok
}
