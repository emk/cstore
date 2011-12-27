package main

import (
	"http"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func assertResponseBody(t *testing.T, expected string, r *http.Response) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		t.Errorf("Can't read stream from %v: %v", r, err)
		return
	}
	assertStringsEqual(t, expected, string(body))
}

func assertHttpGetStatus(t *testing.T, client *http.Client, code int, url string) {
	r, err := client.Get(url)
	if err != nil {
		t.Errorf("Can't GET %v: %v", url, err)
		return
	}
	defer r.Body.Close()
	if r.StatusCode != code {
		t.Errorf("Expected status %v, got %v", code, r.Status)
	}
}

func assertHttpGet(t *testing.T, client *http.Client, expected, url string) {
	r, err := client.Get(url)
	if err != nil {
		t.Errorf("Can't GET %v: %v", url, err)
		return
	}
	defer r.Body.Close()
	if r.StatusCode != http.StatusOK {
		t.Errorf("Unexpected HTTP status: %v", r.Status)
		return
	}
	assertResponseBody(t, expected, r)
}

// Send an HTTP PUT request.  The caller is responsible for calling
// resp.Body.Close().
func put(client *http.Client, url, data string) (resp *http.Response, err os.Error) {
	req, err := http.NewRequest("PUT", url, strings.NewReader(data))
	if err != nil {
		return
	}
	resp, err = client.Do(req)
	return
}

func assertHttpPutStatus(t *testing.T, client *http.Client, code int, url, data string) {
	resp, err := put(client, url, data)
	if err != nil {
		t.Fatalf("Can't PUT data: %s", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != code {
		t.Fatalf("Unexpected HTTP response: %s", resp.Status)
	}
}

func TestServer(t *testing.T) {
	// Create a new server.
	server := NewTestServer()
	defer server.Close()

	// Define our data and where to put it.
	data := "Testing!"
	hash := Digest(data)
	url := server.URL + "/" + hash
	clearRegistryForTest(t, NewRegistry(), hash)

	client := new(http.Client)

	// GET and PUT serveral invalid URLs.
	assertHttpGetStatus(t, client, http.StatusForbidden, server.URL+"/foo")
	assertHttpGetStatus(t, client, http.StatusForbidden, server.URL+"/---")
	assertHttpPutStatus(t, client, http.StatusForbidden, server.URL+"/foo", "data")

	// GET an unknown digest.
	assertHttpGetStatus(t, client, http.StatusNotFound, url)

	// PUT our data to the server.
	assertHttpPutStatus(t, client, http.StatusCreated, url, data)

	// TODO: Test partial writes followed by dropped connections.

	// Make sure it returns something when called.
	assertHttpGet(t, client, data, url)

	// PUT data to the wrong SHA256 sum.
	assertHttpPutStatus(t, client, http.StatusBadRequest, url, "Bogus data")
}

func TestReplication(t *testing.T) {
	// Create two servers.
	server1 := NewTestServer()
	defer server1.Close()
	server2 := NewTestServer()
	defer server2.Close()

	// Define our data and where to put it.
	data := "Testing!"
	hash := Digest(data)
	url1 := server1.URL + "/" + hash
	url2 := server2.URL + "/" + hash
	clearRegistryForTest(t, NewRegistry(), hash)

	client := new(http.Client)

	// Store the data on one server, and read it back on another.
	assertHttpPutStatus(t, client, http.StatusCreated, url1, data)
	assertHttpGet(t, client, data, url2)
}
