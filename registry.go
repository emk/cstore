package main

import (
	"godis"
	"os"
	"sync"
)

type Registry struct {
	hostname string
	locker   sync.Mutex
	client   *godis.Client
}

func NewRegistry(hostname string) *Registry {
	return &Registry{hostname: hostname, client: godis.New("", 0, "")}
}

func (r *Registry) RegisterServer(digest string) (err os.Error) {
	r.locker.Lock()
	defer r.locker.Unlock()

	_, err = r.client.Sadd(digest, r.hostname)
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
