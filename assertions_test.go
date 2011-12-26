package cstore

import (
	"testing"
)

func assertStringsEqual(t *testing.T, expected, got string) {
	if expected != got {
		t.Errorf("Expected %#v, got %#v", expected, got)
	}
}
