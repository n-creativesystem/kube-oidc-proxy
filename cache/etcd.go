package cache

// import (
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"path"
// 	"time"

// 	"github.com/n-creativesystem/oidc-proxy/logger"

// 	"go.etcd.io/etcd/clientv3"
// )

// type etcdCache struct {
// 	cfg    clientv3.Config
// 	minute int64
// 	log    logger.ILogger
// 	prefix string
// }

// type etcdClientFn func(client *clientv3.Client) error

// var _ Cache = &etcdCache{}

// func newEtcdCache(builder *Builder) (Cache, error) {
// 	etcdConfig := clientv3.Config{
// 		Endpoints:   builder.endpoints,
// 		TLS:         builder.cc,
// 		DialTimeout: builder.dialTimeout,
// 	}
// 	if builder.username != "" && builder.password != "" {
// 		etcdConfig.Username = builder.username
// 		etcdConfig.Password = builder.password
// 	}
// 	c := &etcdCache{
// 		cfg:    etcdConfig,
// 		minute: builder.minute,
// 		log:    builder.log,
// 		prefix: builder.prefix,
// 	}
// 	return c, nil
// }

// func (c *etcdCache) new(fn etcdClientFn) error {
// 	client, err := clientv3.New(c.cfg)
// 	if err != nil {
// 		return err
// 	}
// 	defer client.Close()
// 	if err := fn(client); err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (c *etcdCache) Get(ctx context.Context, originalKey string) (string, error) {
// 	var value string
// 	var put bool = false
// 	if err := c.new(func(client *clientv3.Client) error {
// 		key := path.Join(c.prefix, originalKey)
// 		resp, err := client.Get(ctx, key)
// 		if err != nil {
// 			return err
// 		}
// 		for _, kv := range resp.Kvs {
// 			var itm item
// 			if err := json.Unmarshal(kv.Value, &itm); err == nil {
// 				value = itm.Value
// 				expires := time.Unix(itm.Expires, 0)
// 				now := time.Now()
// 				if now.Sub(expires).Minutes() >= 20 {
// 					put = true
// 				}
// 				break
// 			}
// 		}
// 		return nil
// 	}); err != nil {
// 		return "", err
// 	}
// 	if put {
// 		c.Put(ctx, originalKey, value)
// 	}
// 	return value, nil
// }

// func (c *etcdCache) Put(ctx context.Context, originalKey string, value string) error {
// 	err := c.new(func(client *clientv3.Client) error {
// 		key := path.Join(c.prefix, originalKey)
// 		item := newItem(value, c.minute)
// 		if _, err := client.Put(ctx, key, item.ToJson()); err != nil {
// 			return err
// 		}
// 		return nil
// 	})
// 	return err
// }

// func (c *etcdCache) Delete(ctx context.Context, originalKey string) error {
// 	return c.new(func(client *clientv3.Client) error {
// 		key := path.Join(c.prefix, originalKey)
// 		_, err := client.Delete(context.Background(), key)
// 		return err
// 	})
// }

// func (c *etcdCache) Close() error {
// 	return nil
// }

// func (c *etcdCache) Expired(now int64) {
// 	c.new(func(client *clientv3.Client) error {
// 		resp, err := client.Get(context.Background(), c.prefix+"/", clientv3.WithFromKey())
// 		if err != nil {
// 			logger.Log.Error(err)
// 		}
// 		for _, kv := range resp.Kvs {
// 			var itm item
// 			if err := json.Unmarshal(kv.Value, &itm); err == nil {
// 				if itm.Expired(now) {
// 					c.log.Info(fmt.Sprintf("has expired %s", string(kv.Key)))
// 					client.Delete(context.Background(), string(kv.Key))
// 				}
// 			}
// 		}
// 		return nil
// 	})
// }
