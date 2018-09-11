package main

import (
	"context"
	"fmt"
	"github.com/smartwalle/etcd4go"
	"github.com/smartwalle/etcd4go/sample/grpc/hw"
	"go.etcd.io/etcd/clientv3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/resolver"
	"time"
)

func main() {
	var config = clientv3.Config{}
	config.Endpoints = []string{"localhost:2379"}
	resolver.Register(etcd4go.NewResolverWithScheme("etcd", config))

	conn, err := grpc.Dial("etcd://my_service/hw", grpc.WithBalancerName("round_robin"), grpc.WithInsecure())
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	c := hw.NewFirstGRPCClient(conn)

	for {
		rsp, err := c.FirstCall(context.Background(), &hw.FirstRequest{Name: "Yang"})
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println(rsp.Message)
		time.Sleep(time.Second * 1)
	}
}
