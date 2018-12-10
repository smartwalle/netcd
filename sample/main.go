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
	id, key, err := c.Register("my_service", "node_1", "123", 5)
	fmt.Println(id, key, err)

	wg.Wait()
}
