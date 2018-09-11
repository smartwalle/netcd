package etcd4go

import (
	"go.etcd.io/etcd/clientv3"
	"google.golang.org/grpc/resolver"
	"path/filepath"
)

const (
	k_DEFAULT_SCHEME = "etcd"
)

type etcdResolver struct {
	scheme string
	c      *Client
}

func NewResolver(etcd *Client) resolver.Builder {
	return NewResolverWithScheme(k_DEFAULT_SCHEME, etcd)
}

func NewResolverWithScheme(scheme string, c *Client) resolver.Builder {
	return &etcdResolver{scheme: scheme, c: c}
}

func (this *etcdResolver) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOption) (resolver.Resolver, error) {
	var key = filepath.Join(target.Scheme, target.Authority, target.Endpoint)
	watchInfo := this.c.Watch(key, clientv3.WithPrefix())

	watchInfo.Handle(func(eventType, key, path string, value []byte) {
		var paths = watchInfo.GetPaths()
		var addList = make([]resolver.Address, 0, len(paths))
		for _, value := range paths {
			var add = resolver.Address{Addr: string(value)}
			addList = append(addList, add)
		}
		cc.NewAddress(addList)
	})
	return this, nil
}

func (this *etcdResolver) Scheme() string {
	return this.scheme
}

func (this *etcdResolver) ResolveNow(option resolver.ResolveNowOption) {
}

func (this *etcdResolver) Close() {
}

// grpc.Dial("scheme://path")
func (this *Client) RegisterWithScheme(scheme, path, addr string, ttl int64) (err error) {
	return this.Register(scheme, filepath.Join(path, addr), addr, ttl)
}

func (this *Client) UnRegisterWithScheme(scheme, path, addr string) (err error) {
	return this.Revoke(scheme, filepath.Join(path, addr))
}
