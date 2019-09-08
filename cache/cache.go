package cache

import "log"

// Cacher ...
type Cacher interface {
	Set(string, []byte) error
	Get(string) ([]byte, error)
	Del(string) error
	GetStat() Stat
}

//TypeInMemCache ...
const TypeInMemCache = "inMemCache"

//ErrKeyNotFound ...
const ErrKeyNotFound = "ErrKeyNotFound"

// NewCacher ...
func NewCacher(typ string) Cacher {
	var c Cacher
	if typ == TypeInMemCache {
		c = newInMemCache()
	}
	if c == nil {
		panic("Unkonwn cache type" + typ)
	}

	log.Println(typ, "ready to serv")
	return c
}
