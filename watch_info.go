package etcd4go

import (
	"go.etcd.io/etcd/clientv3"
	"sync"
)

const (
	EVENT_TYPE_PUT    = "put"
	EVENT_TYPE_DELETE = "delete"
)

type Handler func(eventType, key, path string, value []byte)

type WatchInfo struct {
	mu      sync.Mutex
	key     string
	paths   map[string][]byte
	handler Handler
	watcher clientv3.Watcher
}

func newWatchInfo(key string, watcher clientv3.Watcher) *WatchInfo {
	var n = &WatchInfo{}
	n.key = key
	n.paths = make(map[string][]byte)
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
		this.handler(EVENT_TYPE_PUT, this.key, path, value)
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
		this.handler(EVENT_TYPE_DELETE, this.key, path, value)
	}
}

func (this *WatchInfo) Handle(h Handler) {
	this.handler = h
	if this.handler != nil {
		for path, value := range this.paths {
			this.handler(EVENT_TYPE_PUT, this.key, path, value)
		}
	}
}

func (this *WatchInfo) Close() error {
	if this.watcher == nil {
		return nil
	}
	return this.watcher.Close()
}
