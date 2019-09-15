package main

import (
	"fmt"
	"gocache/cache"
	"gocache/xhttp"
	"gocache/xtcp"
)

func main() {
	fmt.Println("cache starting ...")
	c := cache.NewCacher(cache.TypeInMemCache)
	go xtcp.NewServer(c).Listen()
	xhttp.NewServer(c).Listen()
}
