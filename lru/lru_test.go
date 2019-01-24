package lru

import (
	"fmt"
	"testing"

	"github.com/t-drk/news_proxy/cache"
)

func TestLRUCache(t *testing.T) {
	var cache cache.Cache = LRU(10)
	for i := 0; i < 10; i++ {
		cache.Add(i, i)
		for j := 0; j <= i; j++ {
			if !cache.Contains(j) {
				t.Errorf(fmt.Sprintf("Cache does not contain %v but expected", j))
			} else {
				fmt.Printf("Cache contains %v as expected.\n", j)
			}
		}
	}
	for i := 0; i < 10; i++ {
		cache.Add(i, i)
		for j := 0; j < 10; j++ {
			if !cache.Contains(j) {
				t.Errorf(fmt.Sprintf("Cache does contains %v but not expected", j))
			} else {
				fmt.Printf("Cache does not contains %v as expected.\n", j)
			}

		}
	}

	for i := 10; i < 20; i++ {
		cache.Add(i, i)
		for j := 0; j < i; j++ {
			if j <= i-10 {
				if cache.Contains(j) {
					t.Errorf(fmt.Sprintf("Cache contains an unexpected element %d", j))
				} else {
					fmt.Printf("Cache does not contais %v as expected.\n", j)
				}

			} else {
				if !cache.Contains(j) {
					t.Errorf(fmt.Sprintf("Cache does not contain %v but expected", j))
				} else {
					fmt.Printf("Cache contains %v as expected.\n", j)
				}

			}
		}
	}
}
