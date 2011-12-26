package cstore

import (
	"fmt"
	"http"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
	"sync"
)

var (
	url_regexp = regexp.MustCompile("^/([0-9a-f]*)$")
)

// Our server state.
type handler struct {
	locker  sync.RWMutex      // Must be held to access content.
	content map[string][]byte // Maps SHA256 digest to content.
}

// Safely store content in our hash table.
func (h *handler) setContent(digest string, content []byte) {
	h.locker.Lock()
	defer h.locker.Unlock()
	h.content[digest] = content
}

// Safely fetch content from our hash table.  Return nil if we don't have
// any content for the specified digest.
func (h *handler) getContent(digest string) []byte {
	h.locker.RLock()
	defer h.locker.RUnlock()
	return h.content[digest]
}

// Read and write content via HTTP.
func (h *handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	log.Println(req.Method, req.URL.Path)

	// Extract the SHA digest from our URL.
	digest, err := parseUrlPath(req.URL.Path)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintln(w, err)
		return
	}

	switch req.Method {
	case "GET":
		h.serveGET(digest, w, req)
	case "PUT":
		h.servePUT(digest, w, req)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// Parse a URL path.  We expect all paths to be "/" followed by an SHA256
// sum.  (We may add support for a bare "/" if we add POST support.)
func parseUrlPath(path string) (digest string, err os.Error) {
	match := url_regexp.FindStringSubmatch(path)
	if match == nil || len(match[1]) != 64 {
		err = os.NewError("Invalid resource path")
		return
	}
	digest = strings.ToLower(match[1])
	return
}

// Attempt to fetch a stored blob.
func (h *handler) serveGET(digest string, w http.ResponseWriter, req *http.Request) {
	content := h.getContent(digest)
	if content == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if _, err := w.Write(content); err != nil {
		log.Println("Error writing response:", err)
		return
	}
}

// Attempt to store a new blob.
func (h *handler) servePUT(digest string, w http.ResponseWriter, req *http.Request) {
	dr := NewDigestReader(req.Body)
	content, err := ioutil.ReadAll(dr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Could not read payload")
		return
	}
	if (digest != dr.Digest()) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "SHA256 digest does not match content!")
		return
	}
	h.setContent(digest, content)
	w.WriteHeader(http.StatusCreated)
}

func NewHandler() http.Handler {
	handler := new(handler)
	handler.content = make(map[string][]byte)
	return handler
}
