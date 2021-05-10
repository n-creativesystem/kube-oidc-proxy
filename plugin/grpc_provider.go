package plugin

import (
	"context"
	"sync"

	"github.com/hashicorp/go-plugin"
	"github.com/n-creativesystem/oidc-proxy/cache"
	proto "github.com/n-creativesystem/oidc-proxy/internal/cache"
	"google.golang.org/grpc"
)

type CacheServer struct {
	Impl cache.Cache
	proto.UnimplementedCacheServer
}

var _ proto.CacheServer = &CacheServer{}

func (c *CacheServer) Init(ctx context.Context, r *proto.SettingRequest) (*proto.Empty, error) {
	c.Impl.Setting(ctx, cache.CacheSetting{
		Endpoints: r.GetEndpoints(),
		Username:  r.GetUserName(),
		Password:  r.GetPassword(),
	})
	return &proto.Empty{}, nil
}

func (c *CacheServer) Get(ctx context.Context, r *proto.GetRequest) (*proto.GetResponse, error) {
	value, err := c.Impl.Get(ctx, r.Key)
	return &proto.GetResponse{
		Value: value,
	}, err
}
func (c *CacheServer) Put(ctx context.Context, r *proto.PutRequest) (*proto.Empty, error) {
	err := c.Impl.Put(ctx, r.Key, r.Value)
	return &proto.Empty{}, err
}
func (c *CacheServer) Delete(ctx context.Context, r *proto.DeleteRequest) (*proto.Empty, error) {
	err := c.Impl.Delete(ctx, r.Key)
	return &proto.Empty{}, err
}

type GRPCCachePlugin struct {
	plugin.Plugin
	GRPCCache func() *CacheServer
}

func (p *GRPCCachePlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &GRPCCache{
		client: proto.NewCacheClient(c),
		ctx:    ctx,
	}, nil
}

func (p *GRPCCachePlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	proto.RegisterCacheServer(s, p.GRPCCache())
	return nil
}

type GRPCCache struct {
	PluginClient *plugin.Client
	TestServer   *grpc.Server
	client       proto.CacheClient
	ctx          context.Context
	mu           sync.Mutex
}

var _ cache.Cache = &GRPCCache{}

func (p *GRPCCache) Get(ctx context.Context, key string) (string, error) {
	r := &proto.GetRequest{
		Key: key,
	}
	res, err := p.client.Get(ctx, r)
	if err != nil {
		return "", err
	}
	return res.Value, nil
}
func (p *GRPCCache) Put(ctx context.Context, key string, value string) error {
	r := &proto.PutRequest{
		Key:   key,
		Value: value,
	}
	_, err := p.client.Put(ctx, r)
	if err != nil {
		return err
	}
	return nil
}
func (p *GRPCCache) Delete(ctx context.Context, key string) error {
	r := &proto.DeleteRequest{
		Key: key,
	}
	_, err := p.client.Delete(ctx, r)
	if err != nil {
		return err
	}
	return nil
}
func (p *GRPCCache) Setting(ctx context.Context, setting cache.CacheSetting) {
	r := &proto.SettingRequest{
		Endpoints: setting.Endpoints,
		UserName:  setting.Username,
		Password:  setting.Password,
		CacheTime: int32(setting.CacheTime),
	}
	p.client.Init(ctx, r)
}

func (p *GRPCCache) Close() error {
	if p.PluginClient == nil {
		return nil
	}
	p.PluginClient.Kill()
	return nil
}
