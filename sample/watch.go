package main

import (
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/smartwalle/etcd4go"
	"sync"
)

func main() {
	var wg = &sync.WaitGroup{}
	wg.Add(1)

	var config = clientv3.Config{}
	config.Endpoints = []string{"localhost:2379"}

	var s, _ = etcd4go.NewService(config)
	info := s.Watch("my_service", clientv3.WithPrefix())

	info.Handle(func(t, key, path string, value []byte) {
		fmt.Println(t, key, path, string(value))
	})

	wg.Wait()
}
