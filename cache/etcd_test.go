package cache_test

import (
	"testing"
)

func TestEtcd(t *testing.T) {
	// buf := new(bytes.Buffer)
	// log := logger.New(buf, logger.Info)
	// c, err := cache.New(cache.Eted, cache.CacheTime(1), cache.Endpoint("http://github.com/n-creativesystem/oidc-proxy-etcd:2380"), cache.Logger(log))
	// assert.NoError(t, err)
	// ctx := context.TODO()
	// err = c.Put(ctx, "test", "test")
	// assert.NoError(t, err)
	// value, err := c.Get(ctx, "test")
	// assert.NoError(t, err)
	// assert.Equal(t, "test", value)
	// time.Sleep(2 * time.Minute)
	// strBuffer := buf.String()
	// buffers := strings.Split(strBuffer, "|")
	// assert.Equal(t, "has expired test\n", strings.Trim(buffers[len(buffers)-1], " "))
	// value, err = c.Get(ctx, "test")
	// assert.NoError(t, err)
	// assert.Equal(t, "", value)
}
