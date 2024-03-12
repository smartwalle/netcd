module github.com/smartwalle/netcd/examples

go 1.12

require (
	github.com/smartwalle/netcd v0.0.0
	go.etcd.io/etcd/client/v3 v3.5.12
)

replace github.com/smartwalle/netcd => ../
