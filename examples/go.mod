module github.com/smartwalle/etcd4go/examples

go 1.12

require (
	go.etcd.io/etcd/client/v3 v3.5.0-alpha.0 // indirect
	github.com/smartwalle/etcd4go v0.0.0
)

replace (
	github.com/smartwalle/etcd4go => ../
)