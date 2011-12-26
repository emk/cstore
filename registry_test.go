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
	digest := Digest("Test.")
	r := NewRegistry()
	assertRegisterServer(t, r, digest, "s1.example.com")
	assertRegisterServer(t, r, digest, "s2.example.com")
	servers, err := r.FindServers(digest)
	if err != nil {
		t.Fatal("Can't find servers:", err)
	}
	assertStringSlicesEqual(t, []string{ "s1.example.com", "s2.example.com" }, servers)
}
