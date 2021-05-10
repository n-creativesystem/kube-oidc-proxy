package cache

// import (
// 	"crypto/tls"
// 	"time"

// 	"github.com/n-creativesystem/oidc-proxy/logger"
// )

// type Builder struct {
// 	endpoints   []string
// 	dialTimeout time.Duration
// 	minute      int64
// 	log         logger.ILogger
// 	cc          *tls.Config
// 	username    string
// 	password    string
// 	prefix      string
// }

// type Option func(builder *Builder)

// func (o Option) Apply(builder *Builder) {
// 	o(builder)
// }

// func Endpoint(endpoints ...string) Option {
// 	return func(builder *Builder) {
// 		builder.endpoints = endpoints
// 	}
// }

// func DialTimeout(timeout time.Duration) Option {
// 	return func(builder *Builder) {
// 		builder.dialTimeout = timeout
// 	}
// }

// func CacheTime(minute int64) Option {
// 	return func(builder *Builder) {
// 		builder.minute = minute
// 	}
// }

// func Logger(log logger.ILogger) Option {
// 	return func(builder *Builder) {
// 		builder.log = log
// 	}
// }

// func TLSConfig(cc *tls.Config) Option {
// 	return func(builder *Builder) {
// 		builder.cc = cc
// 	}
// }

// func Username(username string) Option {
// 	return func(builder *Builder) {
// 		builder.username = username
// 	}
// }

// func Password(password string) Option {
// 	return func(builder *Builder) {
// 		builder.password = password
// 	}
// }

// func Prefix(prefix string) Option {
// 	if prefix == "" {
// 		prefix = defaultPrefix
// 	}
// 	return func(builder *Builder) {
// 		builder.prefix = prefix
// 	}
// }
