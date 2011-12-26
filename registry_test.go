package cstore

import (
	"testing"
)

func assertRegisterServer(t *testing.T, r *Registry, digest, hostname string) {
	if err := r.RegisterServer(digest, hostname); err != nil {
		t.Errorf("Error registering %v: %v", hostname, err)
	}

}

func TestRegistry(t *testing.T) {
	// Register two servers as owning a digest.
	digest := Digest("Test.")
	r := NewRegistry()
	assertRegisterServer(t, r, digest, "s1.example.com")
	assertRegisterServer(t, r, digest, "s2.example.com")

	// Get the complete list of servers.
	servers, err := r.FindServers(digest)
	if err != nil {
		t.Fatal("Can't find servers:", err)
	}
	assertStringSlicesEqual(t, []string{"s1.example.com", "s2.example.com"}, servers)

	// Get a random server holding a specific digest.
	server, err := r.FindOneServer(digest)
	if err != nil {
		t.Fatal("Can't find server:", err)
	}
	if server != "s1.example.com" && server != "s2.example.com" {
		t.Error("Unexpected server value: %#v", server)
	}
}
