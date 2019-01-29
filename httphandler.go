package main

import (
	"net/http"

	"github.com/apex/log"
	"github.com/gorilla/mux"
)

// HTTPHandler - handles incoming HTTP requests
type HTTPHandler struct {
	Router *mux.Router
}

func handleIndex(w http.ResponseWriter, r *http.Request) {

}

// NewHTTPHandler - creates a new HTTPHandler xD
func NewHTTPHandler() *HTTPHandler {
	router := mux.NewRouter()

	router.HandleFunc("/", handleIndex)

	return &HTTPHandler{
		Router: router,
	}
}

// ServeHTTP - http.ListenAndServe invokes this on incoming request
func (h *HTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var err error
	defer log.WithFields(log.Fields{
		"addr": r.RemoteAddr,
		"path": r.URL.Path,
	}).Trace("serving http").Stop(&err)

	h.Router.ServeHTTP(w, r)
}
