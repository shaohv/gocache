package main

import (
	"flag"
	"fmt"
	"gocache/cache"
	"gocache/xhttp"
	"gocache/xtcp"
)

func main() {
	typ := flag.String("type", "inmemory", "cache type")
	flag.Parse()
	fmt.Println("cache starting ...", *typ)
	c := cache.NewCacher(*typ)
	go xtcp.NewServer(c).Listen()
	xhttp.NewServer(c).Listen()
}
