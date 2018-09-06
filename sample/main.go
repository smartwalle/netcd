package main

import (
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/smartwalle/etcd4go"
	"sync"
	"time"
)

func main() {
	var wg = &sync.WaitGroup{}
	wg.Add(1)

	var config = clientv3.Config{}
	config.Endpoints = []string{"localhost:2379"}

	var s, _ = etcd4go.NewService(config)
	fmt.Println(s.Register("my_service", "service1", "123", 5))
	time.Sleep(time.Second * 5)
	//fmt.Println(s.Revoke("my_service", "service1"))
	fmt.Println(s.Register("my_service", "service2", "456", 5))

	wg.Wait()
}
