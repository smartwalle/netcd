package main

import (
	"context"
	"fmt"
	"github.com/smartwalle/netcd"
	"go.etcd.io/etcd/client/v3"
	"os"
	"os/signal"
	"syscall"
	"time"
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

	for i := 0; i < 5; i++ {
		go func(idx int) {
			var lock, merr = client.NewMutex("test-lock")
			if merr != nil {
				fmt.Println("创建互斥锁异常:", merr)
				return
			}

			if merr = lock.Lock(context.Background()); merr != nil {
				fmt.Println("添加锁异常:", merr)
				return
			}

			fmt.Println(idx, " start")
			time.Sleep(time.Second * 4)
			fmt.Println(idx, " end")
			lock.Unlock(context.Background())
		}(i)
	}

	var sigs = make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
MainLoop:
	for {
		s := <-sigs
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			break MainLoop
		}
	}
}
