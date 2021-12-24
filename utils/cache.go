package utils

import (
	"github.com/bluele/gcache"
)

var GCache gcache.Cache

func InitCache() {
	GCache = gcache.New(20).
		LRU().
		Build()
}
