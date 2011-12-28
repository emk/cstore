package main

import (
	"testing"
)

func assertRegisterServer(t *testing.T, r *Registry, digest string) {
	if err := r.RegisterServer(digest); err != nil {
		t.Error("Error registering digest:", err)
	}
}

func TestRegistry(t *testing.T) {
	r1 := NewRegistry("s1.example.com")
	r2 := NewRegistry("s2.example.com")
	clearRegistryForTest(t, r1)

	// Register two servers as owning a digest.
	digest := Digest("Test.")
	unknown_digest := Digest("Unknown")
	assertRegisterServer(t, r1, digest)
	assertRegisterServer(t, r2, digest)

	// Get a random server holding a specific digest.
	server, err := r1.FindOneServer(digest)
	if err != nil {
		t.Fatal("Can't find server:", err)
	}
	if server != "s1.example.com" && server != "s2.example.com" {
		t.Error("Unexpected server value:", server)
	}

	// Ask for a random server when none is registered.
	server, err = r1.FindOneServer(unknown_digest)
	if err != nil {
		t.Error("Unexpected error:", err)
	}
	if server != "" {
		t.Error("Did not expect to find server:", server)
	}
}
