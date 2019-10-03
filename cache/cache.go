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
const (
	TypeInMemCache   = "inmemory"
	TypeRocksdbCache = "rocksdb"
)

//ErrKeyNotFound ...
const ErrKeyNotFound = "ErrKeyNotFound"

// NewCacher ...
func NewCacher(typ string) Cacher {
	var c Cacher
	if typ == TypeInMemCache {
		c = newInMemCache()
	} else if typ == TypeRocksdbCache {
		c = newRocksdbCache()
	}
	if c == nil {
		panic("Unkonwn cache type" + typ)
	}

	log.Println(typ, "ready to serv")
	return c
}
