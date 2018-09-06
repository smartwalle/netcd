package etcd4go

import "sync"

const (
	EventTypePut    = "put"
	EventTypeDelete = "delete"
)

type WatchHandle func(eventType, key, path string, value []byte)

type WatchInfo struct {
	mu     sync.Mutex
	key    string
	paths  map[string][]byte
	handle WatchHandle
}

func newWatchInfo(key string) *WatchInfo {
	var n = &WatchInfo{}
	n.key = key
	n.paths = make(map[string][]byte)
	return n
}

func (this *WatchInfo) Key() string {
	return this.key
}

func (this *WatchInfo) AddPath(path string, value []byte) {
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

func (this *WatchInfo) DeletePath(path string) {
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
}
