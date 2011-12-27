package main

import (
	"godis"
	"os"
	"sync"
)

type Registry struct {
	locker sync.Mutex
	client *godis.Client
}

func NewRegistry() *Registry {
	return &Registry{client: godis.New("", 0, "")}
}

func (r *Registry) RegisterServer(digest, hostname string) (err os.Error) {
	r.locker.Lock()
	defer r.locker.Unlock()

	_, err = r.client.Sadd(digest, hostname)
	return
}

func (r *Registry) FindOneServer(digest string) (server string, err os.Error) {
	r.locker.Lock()
	defer r.locker.Unlock()

	elem, err := r.client.Srandmember(digest)
	if err != nil {
		return
	}
	server = elem.String()
	return
}

func (r *Registry) ClearServers(digest string) (err os.Error) {
	r.locker.Lock()
	defer r.locker.Unlock()

	_, err = r.client.Del(digest)
	return
}
