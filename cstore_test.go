package cstore

import (
	"crypto"
	"encoding/hex"
	"http"
	"http/httptest"
	"io/ioutil"
	"strings"
	"testing"
)

func digest(data string) string {
	hash := crypto.SHA256.New()
	if _, err := hash.Write([]byte(data)); err != nil {
		panic("Writing to a hash should never fail")
	}
	return hex.EncodeToString(hash.Sum())
}

func assertStringsEqual(t *testing.T, expected, got string) {
	if expected != got {
		t.Errorf("Expected %#v, got %#v", expected, got)
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

func assertHttpStatus(t *testing.T, client *http.Client, code int, url string) {
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
	assertResponseBody(t, expected, r)
}

func TestServer(t *testing.T) {
	// Create a new server.
	server := httptest.NewServer(NewHandler())
	defer server.Close()

	// Define our data and where to put it.
	data := "Testing!"
	hash := digest(data)
	url := server.URL + "/" + hash

	client := new(http.Client)

	// GET serveral invalid URLs.
	assertHttpStatus(t, client, http.StatusForbidden, server.URL+"/foo")
	assertHttpStatus(t, client, http.StatusForbidden, server.URL+"/---")

	// GET an unknown digest.
	assertHttpStatus(t, client, http.StatusNotFound, url)

	// PUT our data to the server.
	req, err := http.NewRequest("PUT", url, strings.NewReader(data))
	if err != nil {
		t.Fatalf("Can't build HTTP request: %s", err)
	}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Can't PUT data: %s", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("Unexpected HTTP response: %s", resp.Status)
	}

	// TODO: Test partial writes followed by dropped connections.

	// Make sure it returns something when called.
	assertHttpGet(t, client, data, url)
}
