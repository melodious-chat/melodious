package main

import (
	"net/http"
	"sync"
)

// Melodious - root structure
type Melodious struct {
	Config    *Config
	Database  *Database
	UserConns *sync.Map
}

// NewMelodious - creates a new Melodious instance
func NewMelodious(cfg *Config) *Melodious {
	return &Melodious{
		Config:    cfg,
		Database:  nil,
		UserConns: &sync.Map{},
	}
}

// ConnectToDB - connects to the database
func (mel *Melodious) ConnectToDB() {
	var err error
	mel.Database, err = NewDatabase(mel, mel.Config.DBAddr)
	if err != nil {
		panic(err)
	}
}

// webServerRunner - An internal function used by RunWebServer
func (mel *Melodious) webServerRunner() {
	h := NewHTTPHandler(mel)
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

// PutConnection - ruts a connection to a pool
func (mel *Melodious) PutConnection(username string, connInfo *ConnInfo) {
	m := &sync.Map{}
	m.Store(connInfo, true)
	mm, loaded := mel.UserConns.LoadOrStore(username, m)
	if !loaded {
	} else if m := mm.(*sync.Map); m != nil {
		m.Store(connInfo, true)
		mel.UserConns.Store(username, m)
	} else {
		mel.UserConns.Store(username, m)
	}
}

// RemoveConnection - removes a connection from a pool
func (mel *Melodious) RemoveConnection(username string, connInfo *ConnInfo) {
	m, loaded := mel.UserConns.Load(username)
	if !loaded {
	} else if m := m.(*sync.Map); m != nil {
		m.Delete(connInfo)
	}
}

// IterateOverConnections - iterates over all connections of a given username
func (mel *Melodious) IterateOverConnections(username string, f func(connInfo *ConnInfo)) {
	m, loaded := mel.UserConns.Load(username)
	if !loaded {
	} else if m := m.(*sync.Map); m != nil {
		m.Range(func(key interface{}, value interface{}) bool {
			if connInfo := key.(*ConnInfo); connInfo != nil {
				f(connInfo)
			}
			return true
		})
	}
}
