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

	c.Watch("my_service", func(eventType, key, path string, value []byte) {
		fmt.Println("1", eventType, key, path, string(value))
	}, clientv3.WithPrefix())

	c.Watch("my_service", func(eventType, key, path string, value []byte) {
		fmt.Println("2", eventType, key, path, string(value))
	}, clientv3.WithPrefix())

	c.Watch("my_service", func(eventType, key, path string, value []byte) {
		fmt.Println("3", eventType, key, path, string(value))
	}, clientv3.WithPrefix())

	wg.Wait()
}
