package etcd4go

import (
	"context"
	"go.etcd.io/etcd/clientv3"
	"path/filepath"
	"sync"
)

type Client struct {
	client        *clientv3.Client
	mu            sync.Mutex
	leaseIdList   map[string]clientv3.LeaseID
	watchInfoList map[string]*WatchInfo
}

func NewClient(cfg clientv3.Config) (*Client, error) {
	c, err := clientv3.New(cfg)
	if err != nil {
		return nil, err
	}
	var s = &Client{}
	s.client = c
	s.leaseIdList = make(map[string]clientv3.LeaseID)
	s.watchInfoList = make(map[string]*WatchInfo)
	return s, nil
}

func (this *Client) Register(root, path, value string, ttl int64) (string, error) {
	return this.RegisterWithKey(filepath.Join(root, path), value, ttl)
}

func (this *Client) RegisterWithKey(key, value string, ttl int64) (string, error) {
	this.mu.Lock()
	defer this.mu.Unlock()

	keepAliveRsp, leaseId, err := this.keepAlive(key, value, ttl)
	if err != nil {
		return "", err
	}
	this.leaseIdList[key] = leaseId
	go func() {
		for {
			select {
			case _, ok := <-keepAliveRsp:
				if ok == false {
					this.revoke(leaseId)
					return
				}
			}
		}
	}()
	return key, err
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

func (this *Client) UnRegister(root, path string) (err error) {
	return this.Revoke(root, path)
}

func (this *Client) UnRegisterWithKey(key string) (err error) {
	return this.RevokeWithKey(key)
}

func (this *Client) Revoke(root, path string) (err error) {
	return this.RevokeWithKey(filepath.Join(root, path))
}

func (this *Client) RevokeWithKey(key string) (err error) {
	this.mu.Lock()
	defer this.mu.Unlock()

	if leaseId, ok := this.leaseIdList[key]; ok {
		delete(this.leaseIdList, key)
		return this.revoke(leaseId)
	}
	return nil
}

func (this *Client) revoke(leaseId clientv3.LeaseID) (err error) {
	lease := clientv3.NewLease(this.client)
	_, err = lease.Revoke(context.Background(), leaseId)
	return err
}

func (this *Client) Watch(key string, opts ...clientv3.OpOption) (watchInfo *WatchInfo) {
	watcher := clientv3.NewWatcher(this.client)
	watchChan := watcher.Watch(context.Background(), key, opts...)

	this.mu.Lock()
	defer this.mu.Unlock()
	watchInfo = this.watchInfoList[key]
	if watchInfo == nil {
		watchInfo = newWatchInfo(key)
		this.watchInfoList[key] = watchInfo

		kv := clientv3.NewKV(this.client)
		rsp, _ := kv.Get(context.Background(), key, opts...)
		if rsp != nil {
			for _, k := range rsp.Kvs {
				watchInfo.AddPath(string(k.Key), k.Value)
			}
		}
	}

	go func(info *WatchInfo) {
		for {
			select {
			case wc, ok := <-watchChan:
				if ok == false {
					return
				}
				for _, event := range wc.Events {
					switch event.Type {
					case clientv3.EventTypePut:
						info.AddPath(string(event.Kv.Key), event.Kv.Value)
					case clientv3.EventTypeDelete:
						info.DeletePath(string(event.Kv.Key))
					}
				}
			}
		}
	}(watchInfo)
	return watchInfo
}
