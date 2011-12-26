package cstore

import (
	"godis"
	"os"
	"sort"
)

type Registry struct {
	client *godis.Client
}

func NewRegistry() *Registry {
	return &Registry{client: godis.New("", 0, "")}
}

func (r *Registry) RegisterServer(digest, hostname string) (err os.Error) {
	_, err = r.client.Sadd(digest, hostname)
	return
}

func (r *Registry) FindServers(digest string) (servers []string, err os.Error) {
	reply, err := r.client.Smembers(digest)
	if err != nil {
		return
	}
	servers = reply.StringArray()
	sort.Sort(sort.StringSlice(servers))
	return
}
