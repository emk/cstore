package cstore

import (
	"testing"
	"http"
	"http/httptest"
	"io/ioutil"
)

func assertStringsEqual(t *testing.T, expected, got string) {
	if expected != got {
		t.Errorf("Expected %v, got %v", expected, got)
	}
}

func assertResponseBody(t *testing.T, expected string, r *http.Response) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		t.Errorf("Can't read stream from %v: %v", r, err)
		return
	}
	assertStringsEqual(t, expected, string(body))
}

func assertHttpGet(t *testing.T, client http.Client, expected, url string) {
	r, err := client.Get(url)
	if err != nil {
		t.Errorf("Can't GET %s: %v", url, err)
		return
	}
	defer r.Body.Close()
	assertResponseBody(t, expected, r)
}

func TestServer(t *testing.T) {
	// Create a new server.
	server := httptest.NewServer(NewHandler())
	defer server.Close()

	// Make sure it returns something when called.
	client := new(http.Client)
	assertHttpGet(t, client, "Testing.\n", server.URL)
}
