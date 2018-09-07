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

	var s, _ = etcd4go.NewService(config)
	fmt.Println(s.Register("my_service", "node_1", "123", 5))

	wg.Wait()
}
