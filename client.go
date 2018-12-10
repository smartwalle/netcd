package etcd4go

import (
	"context"
	"go.etcd.io/etcd/clientv3"
	"path/filepath"
	"sync"
)

type Client struct {
	client      *clientv3.Client
	mu          sync.Mutex
	leaseIdList map[clientv3.LeaseID]string
}

func NewClient(cfg clientv3.Config) (*Client, error) {
	c, err := clientv3.New(cfg)
	if err != nil {
		return nil, err
	}
	var s = &Client{}
	s.client = c
	s.leaseIdList = make(map[clientv3.LeaseID]string)
	return s, nil
}

func (this *Client) Register(root, path, value string, ttl int64) (int64, string, error) {
	return this.RegisterWithKey(filepath.Join(root, path), value, ttl)
}

func (this *Client) RegisterWithKey(key, value string, ttl int64) (int64, string, error) {
	this.mu.Lock()
	defer this.mu.Unlock()

	keepAliveRsp, leaseId, err := this.keepAlive(key, value, ttl)
	if err != nil {
		return 0, "", err
	}
	this.leaseIdList[leaseId] = key
	go func(leaseId clientv3.LeaseID, rsp <-chan *clientv3.LeaseKeepAliveResponse) {
		for {
			select {
			case _, ok := <-rsp:
				if ok == false {
					this.Revoke(int64(leaseId))
					return
				}
			}
		}
	}(leaseId, keepAliveRsp)
	return int64(leaseId), key, err
}

func (this *Client) keepAlive(key, value string, ttl int64) (rsp <-chan *clientv3.LeaseKeepAliveResponse, leaseId clientv3.LeaseID, err error) {
	kv := clientv3.NewKV(this.client)
	lease := clientv3.NewLease(this.client)

	grantRsp, err := lease.Grant(context.Background(), ttl)
	if err != nil {
		return nil, 0, err
	}

	if _, err = kv.Put(context.Background(), key, value, clientv3.WithLease(grantRsp.ID)); err != nil {
		return nil, 0, err
	}

	rsp, err = lease.KeepAlive(context.Background(), grantRsp.ID)
	return rsp, grantRsp.ID, err
}

func (this *Client) UnRegister(leaseId int64) (err error) {
	return this.Revoke(leaseId)
}

func (this *Client) Revoke(leaseId int64) (err error) {
	this.mu.Lock()
	defer this.mu.Unlock()

	delete(this.leaseIdList, clientv3.LeaseID(leaseId))

	lease := clientv3.NewLease(this.client)
	_, err = lease.Revoke(context.Background(), clientv3.LeaseID(leaseId))
	return err
}

func (this *Client) Watch(key string, opts ...clientv3.OpOption) (watchInfo *WatchInfo) {
	watcher := clientv3.NewWatcher(this.client)
	watchChan := watcher.Watch(context.Background(), key, opts...)

	this.mu.Lock()
	defer this.mu.Unlock()

	watchInfo = newWatchInfo(key)
	kv := clientv3.NewKV(this.client)
	rsp, _ := kv.Get(context.Background(), key, opts...)
	if rsp != nil {
		for _, k := range rsp.Kvs {
			watchInfo.AddPath(string(k.Key), k.Value)
		}
	}

	var ctx, cancel = context.WithCancel(context.Background())
	watchInfo.cancel = cancel

	go func(ctx context.Context, wi *WatchInfo, wc clientv3.WatchChan) {
		for {
			select {
			case wc, ok := <-wc:
				if ok == false {
					return
				}
				for _, event := range wc.Events {
					switch event.Type {
					case clientv3.EventTypePut:
						wi.AddPath(string(event.Kv.Key), event.Kv.Value)
					case clientv3.EventTypeDelete:
						wi.DeletePath(string(event.Kv.Key))
					}
				}
			case <-ctx.Done():
				return
			}
		}
	}(ctx, watchInfo, watchChan)
	return watchInfo
}
