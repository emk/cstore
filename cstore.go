package cstore

import (
	"http"
	"io"
)

func dummyHandler(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "Testing.\n")
}

func NewHandler() http.Handler {
	return http.HandlerFunc(dummyHandler)
}
