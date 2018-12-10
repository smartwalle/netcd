package etcd4go

import (
	"go.etcd.io/etcd/clientv3"
	"sync"
)

const (
	EventTypePut    = "put"
	EventTypeDelete = "delete"
)

type WatchHandle func(eventType, key, path string, value []byte)

type WatchInfo struct {
	mu      sync.Mutex
	key     string
	paths   map[string][]byte
	handle  WatchHandle
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
	if this.handle != nil {
		this.handle(EventTypePut, this.key, path, value)
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
	if this.handle != nil {
		this.handle(EventTypeDelete, this.key, path, value)
	}
}

func (this *WatchInfo) Handle(h WatchHandle) {
	this.handle = h
	if this.handle != nil {
		for path, value := range this.paths {
			this.handle(EventTypePut, this.key, path, value)
		}
	}
}

func (this *WatchInfo) Close() error {
	if this.watcher == nil {
		return nil
	}
	return this.watcher.Close()
}
