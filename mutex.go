package netcd

import (
	"context"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
)

type Mutex struct {
	mutex   *concurrency.Mutex
	session *concurrency.Session
}

func NewMutex(client *clientv3.Client, key string, opts ...concurrency.SessionOption) (*Mutex, error) {
	var session, err = concurrency.NewSession(client, opts...)
	if err != nil {
		return nil, err
	}
	var mutex = &Mutex{}
	mutex.mutex = concurrency.NewMutex(session, key)
	mutex.session = session
	return mutex, nil
}

func (m *Mutex) TryLock(ctx context.Context) error {
	return m.mutex.TryLock(ctx)
}

func (m *Mutex) Lock(ctx context.Context) error {
	return m.mutex.Lock(ctx)
}

func (m *Mutex) Unlock(ctx context.Context) error {
	var err = m.mutex.Unlock(ctx)
	m.session.Close()
	return err
}
