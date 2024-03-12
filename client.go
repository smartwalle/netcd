package netcd

import (
	"context"
	"go.etcd.io/etcd/client/v3"
)

type Client struct {
	client *clientv3.Client
}

func NewClient(client *clientv3.Client) *Client {
	var c = &Client{}
	c.client = client
	return c
}

func (c *Client) NewKV() clientv3.KV {
	return clientv3.NewKV(c.client)
}

func (c *Client) Register(ctx context.Context, key, value string, ttl int64) (int64, string, error) {
	if ttl <= 0 {
		if _, err := c.client.Put(ctx, key, value); err != nil {
			return 0, "", err
		}
		return 0, key, nil
	}

	keepAliveRsp, leaseId, err := c.keepAlive(ctx, key, value, ttl)
	if err != nil {
		return 0, "", err
	}
	go func() {
		for range keepAliveRsp {
		}
	}()
	return int64(leaseId), key, err
}

func (c *Client) keepAlive(ctx context.Context, key, value string, ttl int64) (rsp <-chan *clientv3.LeaseKeepAliveResponse, leaseId clientv3.LeaseID, err error) {
	grantRsp, err := c.client.Grant(ctx, ttl)
	if err != nil {
		return nil, 0, err
	}

	if _, err = c.client.Put(ctx, key, value, clientv3.WithLease(grantRsp.ID)); err != nil {
		return nil, 0, err
	}

	rsp, err = c.client.KeepAlive(ctx, grantRsp.ID)
	return rsp, grantRsp.ID, err
}

func (c *Client) Unregister(ctx context.Context, key string, opts ...clientv3.OpOption) (err error) {
	_, err = c.client.Delete(ctx, key, opts...)
	return err
}

func (c *Client) Revoke(ctx context.Context, leaseId int64) (err error) {
	_, err = c.client.Revoke(ctx, clientv3.LeaseID(leaseId))
	return err
}

func (c *Client) Watch(ctx context.Context, key string, handler Handler, opts ...clientv3.OpOption) *Watcher {
	var watcher = clientv3.NewWatcher(c.client)
	var watchChan = watcher.Watch(ctx, key, opts...)
	var nWatcher = newWatcher(key, handler, watcher)

	rsp, _ := c.client.Get(ctx, key, opts...)
	if rsp != nil {
		for _, k := range rsp.Kvs {
			nWatcher.add(string(k.Key), k.Value, false)
		}
	}

	go func() {
		for value := range watchChan {
			for _, event := range value.Events {
				switch event.Type {
				case clientv3.EventTypePut:
					nWatcher.add(string(event.Kv.Key), event.Kv.Value, true)
				case clientv3.EventTypeDelete:
					nWatcher.delete(string(event.Kv.Key))
				}
			}
		}
	}()
	return nWatcher
}
