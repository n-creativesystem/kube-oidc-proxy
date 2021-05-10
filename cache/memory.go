package cache

// import (
// 	"context"
// 	"fmt"
// 	"path"
// 	"strings"
// 	"sync"

// 	"github.com/n-creativesystem/oidc-proxy/logger"
// )

// type memoryCache struct {
// 	items  map[string]*item
// 	mu     sync.Mutex
// 	minute int64
// 	log    logger.ILogger
// 	prefix string
// }

// var _ Cache = &memoryCache{}

// func newMemoryCache(builder *Builder) *memoryCache {
// 	c := &memoryCache{
// 		items:  make(map[string]*item),
// 		mu:     sync.Mutex{},
// 		minute: builder.minute,
// 		log:    builder.log,
// 		prefix: builder.prefix,
// 	}
// 	return c
// }

// func (c *memoryCache) Get(ctx context.Context, originalKey string) (string, error) {
// 	c.mu.Lock()
// 	key := path.Join(c.prefix, originalKey)
// 	var s string = ""
// 	if v, ok := c.items[key]; ok {
// 		s = v.Value
// 	}
// 	c.mu.Unlock()
// 	return s, nil
// }

// func (c *memoryCache) Put(ctx context.Context, originalKey string, value string) error {
// 	c.mu.Lock()
// 	key := path.Join(c.prefix, originalKey)
// 	c.items[key] = newItem(value, c.minute)
// 	c.mu.Unlock()
// 	return nil
// }

// func (c *memoryCache) Delete(ctx context.Context, key string) error {
// 	c.mu.Lock()
// 	defer c.mu.Unlock()
// 	key = path.Join(c.prefix, key)
// 	delete(c.items, key)
// 	return nil
// }

// func (c *memoryCache) Close() error { return nil }

// func (c *memoryCache) Expired(now int64) {
// 	c.mu.Lock()
// 	for k, v := range c.items {
// 		if strings.HasPrefix(k, c.prefix) {
// 			if v.Expired(now) {
// 				c.log.Info(fmt.Sprintf("has expired %s", k))
// 				delete(c.items, k)
// 			}
// 		}
// 	}
// 	c.mu.Unlock()
// }
