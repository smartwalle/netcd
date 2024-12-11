package main

import (
	"context"
	"fmt"
	"github.com/smartwalle/netcd"
	"go.etcd.io/etcd/client/v3"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	var config = clientv3.Config{}
	config.Endpoints = []string{"192.168.1.9:2379"}
	etcdClient, err := clientv3.New(config)
	if err != nil {
		fmt.Println(err)
		return
	}

	var client = netcd.NewClient(etcdClient)

	var watcher, _ = client.Watch(context.Background(), "my_service", func(watcher *netcd.Watcher, eventType, key, path string, value []byte) {
		fmt.Println("1", eventType, key, path, string(value))
	}, clientv3.WithPrefix())

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
	watcher.Close()

	etcdClient.Close()
}
