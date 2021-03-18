package etcd4go

import (
	"github.com/coreos/etcd/clientv3"
	"sync"
)

const (
	EventPut    string = "PUT"
	EventDelete string = "DELETE"
)

type Handler func(watcher *Watcher, event, key, path string, value []byte)

type Watcher struct {
	watcher clientv3.Watcher
	key     string
	mu      *sync.RWMutex
	values  map[string][]byte
	handler Handler
}

func newWatcher(key string, handler Handler, watcher clientv3.Watcher) *Watcher {
	var n = &Watcher{}
	n.watcher = watcher
	n.key = key
	n.mu = &sync.RWMutex{}
	n.values = make(map[string][]byte)
	n.handler = handler
	return n
}

func (this *Watcher) Key() string {
	return this.key
}

func (this *Watcher) add(path string, value []byte) {
	this.mu.Lock()
	this.values[path] = value
	this.mu.Unlock()
	if this.handler != nil {
		this.handler(this, EventPut, this.key, path, value)
	}
}

func (this *Watcher) Get(path string) []byte {
	this.mu.RLock()
	defer this.mu.RUnlock()
	return this.values[path]
}

func (this *Watcher) Values() map[string][]byte {
	this.mu.RLock()
	defer this.mu.RUnlock()
	var nMap = make(map[string][]byte)
	for key, value := range this.values {
		nMap[key] = value
	}
	return nMap
}

func (this *Watcher) delete(path string) {
	this.mu.Lock()
	var value = this.values[path]
	delete(this.values, path)
	this.mu.Unlock()
	if this.handler != nil {
		this.handler(this, EventDelete, this.key, path, value)
	}
}

func (this *Watcher) Close() error {
	this.mu.Lock()
	defer this.mu.Unlock()
	if this.watcher == nil {
		return nil
	}
	return this.watcher.Close()
}
