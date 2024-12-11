package main

import (
	"context"
	"fmt"
	"github.com/smartwalle/netcd"
	"go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	var config = clientv3.Config{}
	config.Endpoints = []string{"127.0.0.1:2379"}
	etcdClient, err := clientv3.New(config)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer etcdClient.Close()

	var client = netcd.NewClient(etcdClient)

	var m, mErr = client.NewMutex("testx", concurrency.WithTTL(1))
	if mErr != nil {
		fmt.Println("创建互斥锁异常:", mErr)
		return
	}
	fmt.Println("等待获取资源")
	if mErr = m.Lock(context.Background()); mErr != nil {
		fmt.Println("添加锁异常", mErr)
		return
	}
	defer m.Unlock(context.Background())

	fmt.Println("good")

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
}
