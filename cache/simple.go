package cache

/*
	Map is simplest implementation of the Cache interface.
*/
type cache map[KEY]VALUE

func Simple() *cache {
	c := cache(make(map[KEY]VALUE))
	return &c
}

func (c *cache) Add(key KEY, value VALUE) {
	// add the entry to map
	(*c)[key] = value
}

func (c *cache) Contains(key KEY) bool {
	// check if entry is in the map
	_, ok := (*c)[key]
	return ok
}

func (c *cache) Get(key KEY) (VALUE, bool) {
	// get the entry from the map
	value, ok := (*c)[key]
	return value, ok
}

func (c *cache) Delete(key KEY) {
	delete(*c, key)
}
