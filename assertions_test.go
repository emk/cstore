package cstore

import (
	"testing"
)

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
