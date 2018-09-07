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
	cfg    clientv3.Config
	conn   resolver.ClientConn
}

func NewResolver(cfg clientv3.Config) resolver.Builder {
	return NewResolverWithScheme(k_DEFAULT_SCHEME, cfg)
}

func NewResolverWithScheme(scheme string, cfg clientv3.Config) resolver.Builder {
	return &etcdResolver{scheme: scheme, cfg: cfg}
}

func (this *etcdResolver) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOption) (resolver.Resolver, error) {
	s, err := NewService(this.cfg)
	if err != nil {
		return nil, err
	}
	this.conn = cc

	var key = filepath.Join(target.Scheme, target.Authority, target.Endpoint)
	watchInfo := s.Watch(key, clientv3.WithPrefix())

	watchInfo.Handle(func(eventType, key, path string, value []byte) {
		var paths = watchInfo.GetPaths()
		var addList = make([]resolver.Address, 0, len(paths))
		for _, value := range paths {
			var add = resolver.Address{Addr: string(value)}
			addList = append(addList, add)
		}
		this.conn.NewAddress(addList)
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
func (this *Service) RegisterScheme(scheme, path, addr string, ttl int64) (err error) {
	return this.Register(scheme, filepath.Join(path, addr), addr, ttl)
}

func (this *Service) UnRegisterScheme(scheme, path, addr string) (err error) {
	return this.Revoke(scheme, filepath.Join(path, addr))
}
