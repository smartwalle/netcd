package etcd4go

import (
	"context"
	"go.etcd.io/etcd/client/v3"
)

type Client struct {
	client *clientv3.Client
}

func NewClient(client *clientv3.Client) *Client {
	var s = &Client{}
	s.client = client
	return s
}

func (client *Client) NewKV() clientv3.KV {
	return clientv3.NewKV(client.client)
}

func (client *Client) Register(ctx context.Context, key, value string, ttl int64) (int64, string, error) {
	if ttl <= 0 {
		var kv = client.NewKV()
		if _, err := kv.Put(ctx, key, value); err != nil {
			return 0, "", err
		}
		return 0, key, nil
	}

	keepAliveRsp, leaseId, err := client.keepAlive(ctx, key, value, ttl)
	if err != nil {
		return 0, "", err
	}
	go func(leaseId clientv3.LeaseID, rsp <-chan *clientv3.LeaseKeepAliveResponse) {
		for {
			if _, ok := <-rsp; ok == false {
				return
			}
		}
	}(leaseId, keepAliveRsp)
	return int64(leaseId), key, err
}

func (client *Client) keepAlive(ctx context.Context, key, value string, ttl int64) (rsp <-chan *clientv3.LeaseKeepAliveResponse, leaseId clientv3.LeaseID, err error) {
	var kv = client.NewKV()
	var lease = clientv3.NewLease(client.client)

	grantRsp, err := lease.Grant(ctx, ttl)
	if err != nil {
		return nil, 0, err
	}

	if _, err = kv.Put(ctx, key, value, clientv3.WithLease(grantRsp.ID)); err != nil {
		return nil, 0, err
	}

	rsp, err = lease.KeepAlive(ctx, grantRsp.ID)
	return rsp, grantRsp.ID, err
}

func (client *Client) Unregister(ctx context.Context, key string, opts ...clientv3.OpOption) (err error) {
	var kv = client.NewKV()
	_, err = kv.Delete(ctx, key, opts...)
	return err
}

func (client *Client) Revoke(ctx context.Context, leaseId int64) (err error) {
	var lease = clientv3.NewLease(client.client)
	_, err = lease.Revoke(ctx, clientv3.LeaseID(leaseId))
	return err
}

func (client *Client) Watch(ctx context.Context, key string, handler Handler, opts ...clientv3.OpOption) (watcher *Watcher) {
	var etcdWatcher = clientv3.NewWatcher(client.client)
	var watchChan = etcdWatcher.Watch(ctx, key, opts...)

	watcher = newWatcher(key, handler, etcdWatcher)
	var kv = client.NewKV()
	rsp, _ := kv.Get(ctx, key, opts...)
	if rsp != nil {
		for _, k := range rsp.Kvs {
			watcher.add(string(k.Key), k.Value, false)
		}
	}

	go func(wi *Watcher, wc clientv3.WatchChan) {
		for c := range watchChan {
			for _, event := range c.Events {
				switch event.Type {
				case clientv3.EventTypePut:
					wi.add(string(event.Kv.Key), event.Kv.Value, true)
				case clientv3.EventTypeDelete:
					wi.delete(string(event.Kv.Key))
				}
			}
		}
	}(watcher, watchChan)
	return watcher
}
