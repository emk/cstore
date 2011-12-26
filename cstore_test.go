package cstore

import (
	"http"
	"http/httptest"
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
	server := httptest.NewServer(NewHandler())
	defer server.Close()

	// Define our data and where to put it.
	data := "Testing!"
	hash := Digest(data)
	url := server.URL + "/" + hash

	client := new(http.Client)

	// GET serveral invalid URLs.
	assertHttpGetStatus(t, client, http.StatusForbidden, server.URL+"/foo")
	assertHttpGetStatus(t, client, http.StatusForbidden, server.URL+"/---")

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
