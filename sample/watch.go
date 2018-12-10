package main

import (
	"fmt"
	"github.com/smartwalle/etcd4go"
	"go.etcd.io/etcd/clientv3"
	"sync"
)

func main() {
	var wg = &sync.WaitGroup{}
	wg.Add(1)

	var config = clientv3.Config{}
	config.Endpoints = []string{"localhost:2379"}

	var c, _ = etcd4go.NewClient(config)

	info1 := c.Watch("my_service", clientv3.WithPrefix())
	info1.Handle(func(t, key, path string, value []byte) {
		fmt.Println("1", t, key, path, string(value))
	})

	info2 := c.Watch("my_service", clientv3.WithPrefix())
	info2.Handle(func(t, key, path string, value []byte) {
		fmt.Println("2", t, key, path, string(value))
	})

	info3 := c.Watch("my_service", clientv3.WithPrefix())
	info3.Handle(func(t, key, path string, value []byte) {
		fmt.Println("3", t, key, path, string(value))

		info3.Cancel()
	})

	wg.Wait()
}
