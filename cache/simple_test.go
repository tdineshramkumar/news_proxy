package cache

/*
	Test case for the simple cache
*/
import "testing"

func TestSimpleCache(t *testing.T) {
	var cache Cache = Simple()
	cache.Add("key1", "value1")
	if !cache.Contains("key1") {
		t.Errorf("Cache does not contain key1.\n")
	}
	if value, ok := cache.Get("key1"); !ok || value.(string) != "value1" {
		t.Errorf("Got (%v, %v) for key1.\n", value, ok)
	}
	if cache.Contains("key2") {
		t.Errorf("Cache contains key2.\n")
	}
}
