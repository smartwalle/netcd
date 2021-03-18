package main

import (
	"fmt"
	"github.com/smartwalle/etcd4go"
	"go.etcd.io/etcd/client/v3"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	var config = clientv3.Config{}
	config.Endpoints = []string{"192.168.1.77:2379"}
	etcdClient, err := clientv3.New(config)
	if err != nil {
		fmt.Println(err)
		return
	}

	var client = etcd4go.NewClient(etcdClient)
	id, key, err := client.Register("my_service/node_1", "123", 5)
	fmt.Println(id, key, err)

	var c = make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
MainLoop:
	for {
		s := <-c
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			break MainLoop
		}
	}
	client.Revoke(id)
}
