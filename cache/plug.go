package cache

// import (
// 	context "context"
// )

// type PluginServer struct {
// 	Impl Cache
// 	UnimplementedCacheServer
// }

// var _ CacheServer = &PluginServer{}

// func (g *PluginServer) Get(ctx context.Context, r *GetRequest) (*GetResponse, error) {
// 	value, err := g.Impl.Get(ctx, r.Key)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &GetResponse{
// 		Value: value,
// 	}, nil
// }
// func (g *PluginServer) Put(ctx context.Context, r *PutRequest) (*Empty, error) {
// 	res := &Empty{}
// 	err := g.Impl.Put(ctx, r.Key, r.Value)
// 	return res, err
// }
// func (g *PluginServer) Delete(ctx context.Context, r *DeleteRequest) (*Empty, error) {
// 	res := &Empty{}
// 	err := g.Impl.Delete(ctx, r.Key)
// 	return res, err
// }

// type PluginClient struct {
// 	client CacheClient
// }

// func NewPluginClient(client CacheClient) *PluginClient {
// 	return &PluginClient{
// 		client: client,
// 	}
// }

// var _ Cache = &PluginClient{}

// func (p *PluginClient) Get(ctx context.Context, key string) (string, error) {
// 	r := &GetRequest{Key: key}
// 	res, err := p.client.Get(ctx, r)
// 	if err != nil {
// 		return "", err
// 	}
// 	return res.Value, nil
// }
// func (p *PluginClient) Put(ctx context.Context, key string, value string) error {
// 	r := &PutRequest{
// 		Key:   key,
// 		Value: value,
// 	}
// 	_, err := p.client.Put(ctx, r)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }
// func (p *PluginClient) Delete(ctx context.Context, key string) error {
// 	r := &DeleteRequest{
// 		Key: key,
// 	}
// 	_, err := p.client.Delete(ctx, r)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }
// func (p *PluginClient) Expired(now int64) {}
// func (p *PluginClient) Close() error {
// 	return nil
// }
