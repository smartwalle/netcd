package etcd4go

import (
	"go.etcd.io/etcd/clientv3"
	"sync"
)

const (
	EventTypePut    string = "PUT"
	EventTypeDelete string = "DELETE"
)

type Handler func(eventType, key, path string, value []byte)

type WatchInfo struct {
	mu      sync.Mutex
	key     string
	paths   map[string][]byte
	handler Handler
	watcher clientv3.Watcher
}

func newWatchInfo(key string, handler Handler, watcher clientv3.Watcher) *WatchInfo {
	var n = &WatchInfo{}
	n.key = key
	n.paths = make(map[string][]byte)
	n.handler = handler
	n.watcher = watcher
	return n
}

func (this *WatchInfo) Key() string {
	return this.key
}

func (this *WatchInfo) addPath(path string, value []byte) {
	this.mu.Lock()
	defer this.mu.Unlock()
	this.paths[path] = value
	if this.handler != nil {
		this.handler(EventTypePut, this.key, path, value)
	}
}

func (this *WatchInfo) GetPath(path string) []byte {
	this.mu.Lock()
	defer this.mu.Unlock()
	return this.paths[path]
}

func (this *WatchInfo) GetPaths() map[string][]byte {
	return this.paths
}

func (this *WatchInfo) deletePath(path string) {
	this.mu.Lock()
	defer this.mu.Unlock()
	var value = this.paths[path]
	delete(this.paths, path)
	if this.handler != nil {
		this.handler(EventTypeDelete, this.key, path, value)
	}
}

func (this *WatchInfo) Close() error {
	this.mu.Lock()
	defer this.mu.Unlock()
	if this.watcher == nil {
		return nil
	}
	return this.watcher.Close()
}
