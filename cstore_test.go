package cstore

import (
	"testing"
	"http"
	"http/httptest"
	"io/ioutil"
)

func assertStringsEqual(t *testing.T, expected string, got string) {
	if expected != got {
		t.Errorf("Expected %v, got %v", expected, got)
	}
}

func assertResponseBody(t *testing.T, expected string, r *http.Response) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		t.Fatalf("Can't read stream from %v: %v", r, err)
	}
	assertStringsEqual(t, "Testing.\n", string(body))
}

func TestServer(t *testing.T) {
	// Create a new server.
	server := httptest.NewServer(NewHandler())
	defer server.Close()

	// Make sure it returns something when called.
	client := new(http.Client)
	r, err := client.Get(server.URL)
	if err != nil {
		t.Fatalf("Can't GET %s: %v", server.URL, err)
	}
	defer r.Body.Close()

	// Read the HTTP result and compare it.
	assertResponseBody(t, "Testing.\n", r)
}
