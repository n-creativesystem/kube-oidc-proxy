package cache_test

import (
	"bytes"
	"context"
	"strings"
	"testing"
	"time"

	"github.com/n-creativesystem/oidc-proxy/cache"
	"github.com/n-creativesystem/oidc-proxy/logger"

	"github.com/stretchr/testify/assert"
)

func TestMemory(t *testing.T) {
	buf := new(bytes.Buffer)
	log := logger.New(buf, logger.Info)
	c := cache.New(cache.CacheTime(1), cache.Logger(log))
	// assert.NoError(t, err)
	ctx := context.Background()
	err := c.Put(ctx, "test", "test")
	assert.NoError(t, err)
	value, err := c.Get(ctx, "test")
	assert.NoError(t, err)
	assert.Equal(t, "test", value)
	time.Sleep(2 * time.Minute)
	strBuffer := buf.String()
	buffers := strings.Split(strBuffer, "|")
	assert.Equal(t, "has expired test\n", strings.Trim(buffers[len(buffers)-1], " "))
	value, err = c.Get(ctx, "test")
	assert.NoError(t, err)
	assert.Equal(t, "", value)
}
