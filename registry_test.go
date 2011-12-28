package main

import (
	"testing"
)

func assertNewRegistry(t *testing.T, hostname string) *Registry {
	r := NewRegistry()
	if err := r.RegisterThisServer(hostname); err != nil {
		t.Fatal("Error registering server:", err)
	}
	return r
}

func assertRegisterDigest(t *testing.T, r *Registry, digest string) {
	if err := r.RegisterDigest(digest); err != nil {
		t.Error("Error registering digest:", err)
	}
}

func TestRegistry(t *testing.T) {
	clearRegistryForTest(t)
	r1 := assertNewRegistry(t, "s1.example.com")
	defer r1.UnregisterThisServer()
	r2 := assertNewRegistry(t, "s2.example.com")
	defer r2.UnregisterThisServer()

	// Register two servers as owning a digest.
	digest := Digest("Test.")
	unknown_digest := Digest("Unknown")
	assertRegisterDigest(t, r1, digest)
	assertRegisterDigest(t, r2, digest)

	// Get a random server holding a specific digest.
	servers, err := r1.FindServers(digest)
	if err != nil {
		t.Fatal("Can't find servers:", err)
	}
	assertStringSlicesEqual(t, []string{"s2.example.com", "s1.example.com"},
		servers)

	// Ask for a random server when none is registered.
	servers, err = r1.FindServers(unknown_digest)
	if err != nil {
		t.Error("Unexpected error:", err)
	}
	assertStringSlicesEqual(t, []string{}, servers)
}
