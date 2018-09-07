package main

import (
	"context"
	"fmt"
	"github.com/smartwalle/etcd4go"
	"github.com/smartwalle/ngx/a/hw"
	"go.etcd.io/etcd/clientv3"
	"google.golang.org/grpc"
	"net"
)

var addr = ":5003"

func main() {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Println(err)
		return
	}

	var config = clientv3.Config{}
	config.Endpoints = []string{"localhost:2379"}
	var s, _ = etcd4go.NewService(config)
	fmt.Println(s.RegisterScheme("etcd", "my_service/hw", addr, 5))

	server := grpc.NewServer()
	hw.RegisterFirstGRPCServer(server, &service{})
	server.Serve(listener)
}

type service struct {
}

func (this *service) FirstCall(ctx context.Context, req *hw.FirstRequest) (*hw.FirstResponse, error) {
	return &hw.FirstResponse{Message: fmt.Sprintf("Hello %s, from %s", req.Name, addr)}, nil
}
