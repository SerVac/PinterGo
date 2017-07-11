package cachePool_tests

import (
	"testing"
	"../PinterGo/utils/cache"
	"os"
)

func cachePoolInit(t *testing.T){
	pool := cache.NewCachePool("/.test")
	pool.Put("test_filename.txt", []byte("tetdt_12391291239"))

	_, err := os.OpenFile("./test/test_filename.txt", os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		t.Error(err)
	}
}