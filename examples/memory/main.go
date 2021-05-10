package main

import (
	"context"
	"encoding/json"
	"path"
	"sync"

	"github.com/n-creativesystem/oidc-proxy/cache"
	mPlugin "github.com/n-creativesystem/oidc-proxy/examples/memory/plugin"
	"github.com/n-creativesystem/oidc-proxy/plugin"
)

type item struct {
	Value   string `json:"value"`
	Expires int64  `json:"expires"`
}

func newItem(value string) *item {
	return &item{
		Value: value,
	}
}

func (i *item) ToJson() string {
	buf, _ := json.Marshal(i)
	return string(buf)
}

type memoryCache struct {
	items  map[string]*item
	mu     sync.Mutex
	prefix string
}

var _ cache.Cache = &memoryCache{}

func (c *memoryCache) Get(ctx context.Context, originalKey string) (string, error) {
	c.mu.Lock()
	key := path.Join(c.prefix, originalKey)
	var s string = ""
	if v, ok := c.items[key]; ok {
		s = v.Value
	}
	c.mu.Unlock()
	return s, nil
}
func (c *memoryCache) Put(ctx context.Context, originalKey string, value string) error {
	c.mu.Lock()
	key := path.Join(c.prefix, originalKey)
	c.items[key] = newItem(value)
	c.mu.Unlock()
	return nil
}
func (c *memoryCache) Delete(ctx context.Context, key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	key = path.Join(c.prefix, key)
	delete(c.items, key)
	return nil
}
func (c *memoryCache) Close() error                                            { return nil }
func (c *memoryCache) Setting(ctx context.Context, setting cache.CacheSetting) {}

func newMemoryCache() *memoryCache {
	c := &memoryCache{
		items:  make(map[string]*item),
		mu:     sync.Mutex{},
		prefix: "memory",
	}
	return c
}

func main() {
	mPlugin.Sever(&mPlugin.ServerOpts{
		GRPCCacheFunc: func() *plugin.CacheServer {
			return &plugin.CacheServer{
				Impl: newMemoryCache(),
			}
		},
	})
}
