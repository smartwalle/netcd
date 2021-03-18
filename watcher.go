package etcd4go

import (
	"github.com/coreos/etcd/clientv3"
	"sync"
)

const (
	EventPut    string = "PUT"
	EventDelete string = "DELETE"
)

type Handler func(event, key, path string, value []byte)

type Watcher struct {
	key     string
	mu      *sync.RWMutex
	values  map[string][]byte
	handler Handler
	watcher clientv3.Watcher
}

func newWatcher(key string, handler Handler, watcher clientv3.Watcher) *Watcher {
	var n = &Watcher{}
	n.key = key
	n.mu = &sync.RWMutex{}
	n.values = make(map[string][]byte)
	n.handler = handler
	n.watcher = watcher
	return n
}

func (this *Watcher) Key() string {
	return this.key
}

func (this *Watcher) add(path string, value []byte) {
	this.mu.Lock()
	defer this.mu.Unlock()
	this.values[path] = value
	if this.handler != nil {
		this.handler(EventPut, this.key, path, value)
	}
}

func (this *Watcher) Get(path string) []byte {
	this.mu.RLock()
	defer this.mu.RUnlock()
	return this.values[path]
}

func (this *Watcher) Values() map[string][]byte {
	this.mu.RLock()
	defer this.mu.Unlock()
	var nMap = make(map[string][]byte)
	for key, value := range this.values {
		nMap[key] = value
	}
	return nMap
}

func (this *Watcher) delete(path string) {
	this.mu.Lock()
	defer this.mu.Unlock()
	var value = this.values[path]
	delete(this.values, path)
	if this.handler != nil {
		this.handler(EventDelete, this.key, path, value)
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
