package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Melodious - root structure
type Melodious struct {
	Config *Config
}

// NewMelodious - creates a new Melodious instance
func NewMelodious(cfg *Config) *Melodious {
	return &Melodious{
		Config: cfg,
	}
}

// webServerRunner - An internal function used by RunWebServer
func (mel *Melodious) webServerRunner() {
	h := &HTTPHandler{Router: mux.NewRouter()}
	http.ListenAndServe(mel.Config.HTTPAddr, h)
}

// RunWebServer - starts an HTTP server
func (mel *Melodious) RunWebServer() <-chan bool {
	done := make(chan bool)

	go func() {
		mel.webServerRunner()
		done <- true
	}()

	return done
}
