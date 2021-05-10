package cache

import (
	"context"
)

// const (
// 	defaultPrefix = "/github.com/n-creativesystem/oidc-proxy"
// )

type CacheSetting struct {
	Endpoints []string
	Username  string
	Password  string
	CacheTime int
}

type Cache interface {
	Get(ctx context.Context, key string) (string, error)
	Put(ctx context.Context, key string, value string) error
	Delete(ctx context.Context, key string) error
	Setting(ctx context.Context, setting CacheSetting)
	Close() error
}

// type CacheType int

// const (
// 	Memory CacheType = iota
// 	Eted
// 	Redis
// )

// func New(opts ...Option) Cache {
// 	var cacheTime int64 = 300
// 	// var c Cache
// 	builder := &Builder{
// 		minute: cacheTime,
// 		prefix: defaultPrefix,
// 	}
// 	for _, opt := range opts {
// 		opt.Apply(builder)
// 	}
// 	return newMemoryCache(builder)
// 	// switch cacheType {
// 	// case Eted:
// 	// 	return nil, nil
// 	// 	// var err error
// 	// 	// // c, err = newEtcdCache(builder)
// 	// 	// if err != nil {
// 	// 	// 	return nil, err
// 	// 	// }
// 	// case Redis:
// 	// 	return nil, nil
// 	// default:
// 	// 	c = newMemoryCache(builder)
// 	// }
// 	// return c, nil
// }
