package etcd4go

import (
	"context"
	"github.com/coreos/etcd/clientv3"
)

type Client struct {
	client *clientv3.Client
}

func NewClient(client *clientv3.Client) *Client {
	var s = &Client{}
	s.client = client
	return s
}

func (this *Client) NewKV() clientv3.KV {
	return clientv3.NewKV(this.client)
}

func (this *Client) Register(key, value string, ttl int64) (int64, string, error) {
	if ttl <= 0 {
		var kv = this.NewKV()
		if _, err := kv.Put(context.Background(), key, value); err != nil {
			return 0, "", err
		}
		return 0, key, nil
	}

	_, leaseId, err := this.keepAlive(key, value, ttl)
	if err != nil {
		return 0, "", err
	}
	return int64(leaseId), key, err
}

func (this *Client) keepAlive(key, value string, ttl int64) (rsp <-chan *clientv3.LeaseKeepAliveResponse, leaseId clientv3.LeaseID, err error) {
	var kv = this.NewKV()
	var lease = clientv3.NewLease(this.client)

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

func (this *Client) Deregister(key string, opts ...clientv3.OpOption) (err error) {
	var kv = this.NewKV()
	_, err = kv.Delete(context.Background(), key, opts...)
	return err
}

func (this *Client) Revoke(leaseId int64) (err error) {
	var lease = clientv3.NewLease(this.client)
	_, err = lease.Revoke(context.Background(), clientv3.LeaseID(leaseId))
	return err
}

func (this *Client) Watch(key string, handler Handler, opts ...clientv3.OpOption) (watcher *Watcher) {
	var etcdWatcher = clientv3.NewWatcher(this.client)
	var watchChan = etcdWatcher.Watch(context.Background(), key, opts...)

	watcher = newWatcher(key, handler, etcdWatcher)
	var kv = this.NewKV()
	rsp, _ := kv.Get(context.Background(), key, opts...)
	if rsp != nil {
		for _, k := range rsp.Kvs {
			watcher.add(string(k.Key), k.Value)
		}
	}

	go func(wi *Watcher, wc clientv3.WatchChan) {
		for c := range watchChan {
			for _, event := range c.Events {
				switch event.Type {
				case clientv3.EventTypePut:
					wi.add(string(event.Kv.Key), event.Kv.Value)
				case clientv3.EventTypeDelete:
					wi.delete(string(event.Kv.Key))
				}
			}
		}
	}(watcher, watchChan)
	return watcher
}
