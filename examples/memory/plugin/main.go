package plugin

import (
	hplugin "github.com/hashicorp/go-plugin"
	"github.com/n-creativesystem/oidc-proxy/plugin"
	"google.golang.org/grpc"
)

const CachePluginName = "cache"

var Handshake = hplugin.HandshakeConfig{
	MagicCookieKey:   "CACHE_PLUGIN",
	MagicCookieValue: "m9erzlkcuac9gy4a2szc19j7xjleo4s4epwiio9opv8tjv9sid0qetl7cjo6ulkiskorqyg26pcsfyf979pgn28s5a7byfbq0n66",
}

type GRPCCacheFunc func() *plugin.CacheServer

type ServerOpts struct {
	GRPCCacheFunc GRPCCacheFunc
	TestConfig    *hplugin.ServeTestConfig
}

func Sever(opts *ServerOpts) {
	provider := opts.GRPCCacheFunc()
	hplugin.Serve(&hplugin.ServeConfig{
		HandshakeConfig: Handshake,
		VersionedPlugins: map[int]hplugin.PluginSet{
			1: {
				CachePluginName: &plugin.GRPCCachePlugin{
					GRPCCache: func() *plugin.CacheServer {
						return provider
					},
				},
			},
		},
		GRPCServer: func(opts []grpc.ServerOption) *grpc.Server {
			return grpc.NewServer(opts...)
		},
		Test: opts.TestConfig,
	})
}
