package main

import (
	"fmt"
	"godis"
	"log"
	"os"
	"sync"
	"time"
)

// An interface to global data shared by all servers.  For now, we store
// all the data in Redis.
type Registry struct {
	locker sync.Mutex    // Must be held to use 'client'.
	client *godis.Client // Our connection to Redis.

	hostname string    // Address used to access this server.
	serverId int64     // Registered ID for this server.
	quit     chan bool // Used to stop heartbeat.
}

// Create a new Registry client.
func NewRegistry() *Registry {
	return &Registry{client: godis.New("", 0, ""), quit: make(chan bool)}
}

// Assign this server a unique ID and register it using 'hostname'.  This
// starts a background process which sends a periodic heartbeat to Redis.
// To stop the heartbeat, call UnregisterThisServer().
func (r *Registry) RegisterThisServer(hostname string) (err os.Error) {
	r.locker.Lock()
	defer r.locker.Unlock()

	if r.serverId > 0 {
		return os.NewError("Server is already registered")
	}

	// Allocate an ID for this server.
	serverId, err := r.client.Incr("cstore:server_id_generator")
	if err != nil {
		return
	}

	// Register our server's hostname.
	err = r.client.Setex(serverKey(serverId), 20, hostname)
	if err != nil {
		return
	}
	r.hostname, r.serverId = hostname, serverId

	// Start our background heartbeat goroutine.
	go r.sendHeartbeats()
	return
}

func serverKey(serverId interface{}) string {
	return fmt.Sprintf("cstore:server:%v", serverId)
}

func (r *Registry) sendHeartbeats() {
	ticker := time.NewTicker(5 * 1e9)
	for {
		select {
		case <-ticker.C:
			if err := r.sendHeartbeat(); err != nil {
				log.Println("Can't send heartbeat:", err)
			}
		case <-r.quit:
			ticker.Stop()
			return
		}
	}
}

func (r *Registry) sendHeartbeat() (err os.Error) {
	r.locker.Lock()
	defer r.locker.Unlock()
	return r.client.Setex(serverKey(r.serverId), 20, r.hostname)
}

// Unregister a server from Redis, and stop the heartbeat.
func (r *Registry) UnregisterThisServer() {
	r.locker.Lock()
	defer r.locker.Unlock()

	if r.serverId == 0 {
		panic("Can't unregister server that was never registered")
	}

	r.quit <- true
	if _, err := r.client.Del(serverKey(r.serverId)); err != nil {
		log.Println("Unable to delete server key:", err)
	}
	r.hostname, r.serverId = "", 0
}

// Register this server as having the specified digest.  You must call
// RegisterServer before calling this function.
func (r *Registry) RegisterDigest(digest string) (err os.Error) {
	r.locker.Lock()
	defer r.locker.Unlock()

	if r.serverId == 0 {
		return os.NewError("Must register server first")
	}

	_, err = r.client.Sadd("cstore:blob:"+digest, r.serverId)
	return
}

// Return a list of servers which _should_ have the specified digest.
func (r *Registry) FindServers(digest string) (servers []string, err os.Error) {
	r.locker.Lock()
	defer r.locker.Unlock()

	// Look up the IDs of the servers which have this digest.
	reply, err := r.client.Smembers("cstore:blob:" + digest)
	if err != nil {
		return
	}
	ids := reply.StringArray()

	// We can't call Mget unless we have at least one key.
	if len(ids) == 0 {
		servers = make([]string, 0)
		return
	}

	// Build a list of server keys.
	serverKeys := make([]string, 0, len(ids))
	for _, id := range ids {
		serverKeys = append(serverKeys, serverKey(id))
	}

	// Look up our server addresses.
	reply, err = r.client.Mget(serverKeys...)
	if err != nil {
		return
	}
	rawServers := reply.StringArray()

	// Filter out blank addresses, which presumably belong to dead
	// servers.
	servers = make([]string, 0, len(rawServers))
	for _, server := range rawServers {
		if server != "" {
			servers = append(servers, server)
		}
	}
	return
}

// Erase everything is Redis.  For testing only.
func (r *Registry) ClearForTest() (err os.Error) {
	r.locker.Lock()
	defer r.locker.Unlock()

	keys, err := r.client.Keys("cstore:*")
	if err != nil {
		return
	}
	for _, k := range keys {
		_, err = r.client.Del(k)
		if err != nil {
			return
		}
	}
	return
}
