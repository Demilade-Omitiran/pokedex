package pokecache

import "time"

type Cache struct {
	createdAt time.Time
	val       []byte
}

func (c *Cache) Add(key string, val []byte) {

}

func (c *Cache) Get(key string) ([]byte, bool) {
	return nil, false
}

func NewCache() *Cache {
	return &Cache{
		createdAt: time.Now(),
	}
}
