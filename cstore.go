package cstore

import (
	"http"
	"io"
	"io/ioutil"
	"log"
	"regexp"
	"strings"
	"sync"
)

var (
	url_regexp = regexp.MustCompile("^/([0-9a-f]*)$")
)

type handler struct {
	content map[string][]byte
	locker  sync.Locker
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
	h.locker.Lock()
	defer h.locker.Unlock()
	return h.content[digest]
}

func (h *handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	log.Println(req.Method, req.URL.Path)
	match := url_regexp.FindStringSubmatch(req.URL.Path)
	if match == nil {
		w.WriteHeader(http.StatusForbidden)
		io.WriteString(w, "Invalid resource path\n")
		return
	}
	digest := strings.ToLower(match[1])
	if len(digest) != 64 {
		w.WriteHeader(http.StatusForbidden)
		io.WriteString(w, "Invalid resource path\n")
		return
	}

	switch req.Method {
	case "GET":
		content := h.getContent(digest)
		if content == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if _, err := w.Write(content); err != nil {
			log.Println("Error writing response:", err)
			return
		}
	case "PUT":
		content, err := ioutil.ReadAll(req.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			io.WriteString(w, "Could not read payload\n")
			return
		}
		// TODO: Check digest.
		h.setContent(digest, content)
		w.WriteHeader(http.StatusCreated)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func NewHandler() http.Handler {
	handler := new(handler)
	handler.locker = new(sync.Mutex)
	handler.content = make(map[string][]byte)
	return handler
}
