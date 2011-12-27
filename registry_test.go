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
	unknown_digest := Digest("Unknown")
	r := NewRegistry()
	clearRegistryForTest(t, r, digest)

	// Register two servers as owning a digest.
	assertRegisterServer(t, r, digest, "s1.example.com")
	assertRegisterServer(t, r, digest, "s2.example.com")

	// Get a random server holding a specific digest.
	server, err := r.FindOneServer(digest)
	if err != nil {
		t.Fatal("Can't find server:", err)
	}
	if server != "s1.example.com" && server != "s2.example.com" {
		t.Error("Unexpected server value:", server)
	}

	// Ask for a random server when none is registered.
	server, err = r.FindOneServer(unknown_digest)
	if err != nil {
		t.Error("Unexpected error:", err)
	}
	if server != "" {
		t.Error("Did not expect to find server:", server)
	}
}
