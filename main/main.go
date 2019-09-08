package main

import (
	"fmt"
	"gocache/cache"
	"gocache/xhttp"
)

func main() {
	fmt.Println("cache starting ...")
	c := cache.NewCacher(cache.TypeInMemCache)
	xhttp.NewServer(c).Listen()
}
