package main

import (
	"testing"
)

func clearRegistryForTest(t *testing.T, r *Registry, digest string) {
	if err := r.ClearServers(digest); err != nil {
		t.Fatal("Can't clear Redis keys:", err)
	}
}

func assertStringsEqual(t *testing.T, expected, got string) {
	if expected != got {
		t.Errorf("Expected %#v, got %#v", expected, got)
	}
}

func assertStringSlicesEqual(t *testing.T, expected, got []string) {
	match := false
	if len(expected) == len(got) {
		for i := range expected {
			if expected[i] != got[i] {
				break
			}
		}
		match = true
	}
	if !match {
		t.Errorf("Expected %#v, got %#v", expected, got)
	}
}
