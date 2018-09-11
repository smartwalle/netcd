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
	info := c.Watch("my_service", clientv3.WithPrefix())

	info.Handle(func(t, key, path string, value []byte) {
		fmt.Println(t, key, path, string(value))
	})

	wg.Wait()
}
