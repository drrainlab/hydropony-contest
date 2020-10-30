package cacher

import (
	"github.com/bradfitz/gomemcache/memcache"
)

// main memcache instance
var Instance *memcache.Client

// InitCache - Initialize cache
func InitCache(connURL string) {
	Instance = memcache.New(connURL)
}
