package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"sync"

	mPlugin "github.com/n-creativesystem/oidc-proxy/examples/memory/plugin"
	"github.com/n-creativesystem/oidc-proxy/logger"
	"github.com/n-creativesystem/oidc-proxy/session"
)

var log logger.ILogger
var file *os.File

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

type memorySession struct {
	items  map[string]*item
	mu     sync.Mutex
	prefix string
}

var _ session.Session = &memorySession{}

func (c *memorySession) Get(ctx context.Context, originalKey string) (string, error) {
	c.mu.Lock()
	key := path.Join(c.prefix, originalKey)
	var s string = ""
	if v, ok := c.items[key]; ok {
		s = v.Value
	}
	log.Debug(fmt.Sprintf("[GET] %s:%s", key, s))
	c.mu.Unlock()
	return s, nil
}
func (c *memorySession) Put(ctx context.Context, originalKey string, value string) error {
	c.mu.Lock()
	key := path.Join(c.prefix, originalKey)
	c.items[key] = newItem(value)
	log.Debug(fmt.Sprintf("[PUT] %s:%s", key, value))
	c.mu.Unlock()
	return nil
}
func (c *memorySession) Delete(ctx context.Context, key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	key = path.Join(c.prefix, key)
	log.Debug(fmt.Sprintf("[DEL] %s", key))
	delete(c.items, key)
	return nil
}
func (c *memorySession) Close(ctx context.Context) error {
	if file != nil {
		log.Debug("close method: file close")
		file.Close()
	}
	return nil
}
func (c *memorySession) Init(ctx context.Context, setting map[string]interface{}) error {
	if prefix, ok := setting["prefix"].(string); ok {
		c.prefix = prefix
	}
	var filename string
	var ok bool
	if filename, ok = setting["filename"].(string); !ok {
		filename = ""
	}
	var write io.Writer
	if filename != "" {
		var err error
		if file, err = os.Create(filename); err != nil {
			write = os.Stdout
		} else {
			write = file
		}
	} else {
		write = os.Stdout
	}
	var logLevel string
	if logLevel, ok = setting["loglevel"].(string); !ok {
		logLevel = logger.Info.String()
	}
	log = logger.New(write, logger.Convert(logLevel), logger.FormatLong, logger.FormatDatetime)
	log.Info(fmt.Sprintf("%#v", setting))
	return nil
}

func newMemorySession() *memorySession {
	c := &memorySession{
		items:  make(map[string]*item),
		mu:     sync.Mutex{},
		prefix: "memory",
	}
	return c
}

func main() {
	mPlugin.Sever(&mPlugin.ServerOpts{
		GRPCSessionFunc: func() session.Session {
			return newMemorySession()
		},
	})
}
