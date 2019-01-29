package main

import (
	"fmt"
	"net/http"

	"github.com/apex/log"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

// HTTPHandler - handles incoming HTTP requests
type HTTPHandler struct {
	Router *mux.Router
}

// handleIndex - Handles clients which want to receive the index page
func handleIndex(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World! You should really use a proper Melodious client instead of opening this page\n")
}

// handleConnect - Handles clients which want to connect to Melodious
func handleConnect(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.WithFields(log.Fields{"err": err, "addr": r.RemoteAddr, "path": r.URL.Path}).Error("cannot upgrade to websocket")
	} else {
		go handleConnection(conn)
	}
}

// NewHTTPHandler - creates a new HTTPHandler xD
func NewHTTPHandler() *HTTPHandler {
	router := mux.NewRouter()

	router.HandleFunc("/", handleIndex)
	router.HandleFunc("/connect", handleConnect)

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
