package etcd4go

import (
	"context"
	"go.etcd.io/etcd/clientv3"
)

type Client struct {
	client *clientv3.Client
}

func NewClient(cfg clientv3.Config) (*Client, error) {
	c, err := clientv3.New(cfg)
	if err != nil {
		return nil, err
	}
	var s = &Client{}
	s.client = c
	return s, nil
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
	//go func(leaseId clientv3.LeaseID, rsp <-chan *clientv3.LeaseKeepAliveResponse) {
	//	for {
	//		select {
	//		case _, ok := <-rsp:
	//			if ok == false {
	//				this.Revoke(int64(leaseId))
	//				return
	//			}
	//		}
	//	}
	//}(leaseId, keepAliveRsp)
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

func (this *Client) Deregister(leaseId int64) (err error) {
	return this.Revoke(leaseId)
}

func (this *Client) Revoke(leaseId int64) (err error) {
	lease := clientv3.NewLease(this.client)
	_, err = lease.Revoke(context.Background(), clientv3.LeaseID(leaseId))
	return err
}

func (this *Client) Watch(key string, handler Handler, opts ...clientv3.OpOption) (watchInfo *WatchInfo) {
	var watcher = clientv3.NewWatcher(this.client)
	var watchChan = watcher.Watch(context.Background(), key, opts...)

	watchInfo = newWatchInfo(key, handler, watcher)
	var kv = this.NewKV()
	rsp, _ := kv.Get(context.Background(), key, opts...)
	if rsp != nil {
		for _, k := range rsp.Kvs {
			watchInfo.addPath(string(k.Key), k.Value)
		}
	}

	go func(wi *WatchInfo, wc clientv3.WatchChan) {
		for {
			select {
			case wc, ok := <-wc:
				if ok == false {
					return
				}
				for _, event := range wc.Events {
					switch event.Type {
					case clientv3.EventTypePut:
						wi.addPath(string(event.Kv.Key), event.Kv.Value)
					case clientv3.EventTypeDelete:
						wi.deletePath(string(event.Kv.Key))
					}
				}
			}
		}
	}(watchInfo, watchChan)
	return watchInfo
}
