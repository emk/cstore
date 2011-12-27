package main

import (
	"io/ioutil"
	"strings"
	"testing"
)

func TestDigest(t *testing.T) {
	d := "314ad142957febe390cc7223b4deb1d1b21c187f84f6e7257a23fe46c27fcae3"
	assertStringsEqual(t, d, Digest("Test."))
}

func TestHashingReader(t *testing.T) {
	r := NewDigestReader(strings.NewReader("Test."))
	out, err := ioutil.ReadAll(r)
	if err != nil {
		t.Fatal("Error reading with NewHashingReader")
	}
	assertStringsEqual(t, "Test.", string(out))
	assertStringsEqual(t, Digest("Test."), r.Digest())
}
