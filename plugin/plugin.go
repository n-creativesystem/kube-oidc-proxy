package plugin

import "github.com/hashicorp/go-plugin"

var VersionedPlugins = map[int]plugin.PluginSet{
	1: {
		"cache": &GRPCCachePlugin{},
	},
}

var Handshake = plugin.HandshakeConfig{
	MagicCookieKey:   "CACHE_PLUGIN",
	MagicCookieValue: "m9erzlkcuac9gy4a2szc19j7xjleo4s4epwiio9opv8tjv9sid0qetl7cjo6ulkiskorqyg26pcsfyf979pgn28s5a7byfbq0n66",
}
