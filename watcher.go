package etcd4go

import (
	"go.etcd.io/etcd/client/v3"
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

func (w *Watcher) Key() string {
	return w.key
}

func (w *Watcher) add(path string, value []byte, notify bool) {
	w.mu.Lock()
	w.values[path] = value
	w.mu.Unlock()
	if notify && w.handler != nil {
		w.handler(w, EventPut, w.key, path, value)
	}
}

func (w *Watcher) Get(path string) []byte {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.values[path]
}

func (w *Watcher) Values() map[string][]byte {
	w.mu.RLock()
	defer w.mu.RUnlock()
	var nMap = make(map[string][]byte)
	for key, value := range w.values {
		nMap[key] = value
	}
	return nMap
}

func (w *Watcher) delete(path string) {
	w.mu.Lock()
	var value = w.values[path]
	delete(w.values, path)
	w.mu.Unlock()
	if w.handler != nil {
		w.handler(w, EventDelete, w.key, path, value)
	}
}

func (w *Watcher) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.watcher == nil {
		return nil
	}
	return w.watcher.Close()
}
