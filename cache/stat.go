package cache

// Stat stores cacher status
type Stat struct {
	Count     int64
	KeySize   int64
	ValueSize int64
}

func (s *Stat) add(key string, value []byte) {
	s.Count++
	s.ValueSize += int64(len(value))
	s.KeySize += int64(len(key))
}

func (s *Stat) del(key string, value []byte) {
	s.Count--
	s.ValueSize -= int64(len(value))
	s.KeySize -= int64(len(key))
}
