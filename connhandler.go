package main

import (
	"fmt"
	"runtime/debug"
	"sync"

	"github.com/apex/log"
	"github.com/gorilla/websocket"
)

// ConnInfo - stores some info about a connection
type ConnInfo struct {
	mel           *Melodious
	connection    *websocket.Conn
	messageStream chan<- BaseMessage
	subscriptions *sync.Map
	loggedIn      bool
	username      string
}

// HasFlag - checks if the given connection has the given flag
func (connInfo *ConnInfo) HasFlag(channel string, flag string) (bool, error) {
	if !connInfo.loggedIn {
		return false, nil
	}
	has, err := connInfo.mel.Database.HasFlag(connInfo.username, channel, flag)
	return has, err
}

// HasPerm - checks if the given connection has the given permission
func (connInfo *ConnInfo) HasPerm(channel string, flag string) (bool, error) {
	if !connInfo.loggedIn {
		return false, nil
	}
	owner, err := connInfo.mel.Database.IsUserOwner(connInfo.username)
	if err != nil {
		return false, err
	}
	if owner {
		return true, nil
	}
	has, err := connInfo.mel.Database.HasFlag(connInfo.username, channel, flag)
	return has, err
}

// HasPermChID - checks if the given connection has the given permission
func (connInfo *ConnInfo) HasPermChID(channel int, flag string) (bool, error) {
	if !connInfo.loggedIn {
		return false, nil
	}
	owner, err := connInfo.mel.Database.IsUserOwner(connInfo.username)
	if err != nil {
		return false, err
	}
	if owner {
		return true, nil
	}
	has, err := connInfo.mel.Database.HasFlagChID(connInfo.username, channel, flag)
	return has, err
}

// handleConnection Handles users which are connected to Melodious
func handleConnection(mel *Melodious, conn *websocket.Conn) {

	messageStream := make(chan BaseMessage)

	connInfo := &ConnInfo{
		mel:           mel,
		connection:    conn,
		messageStream: messageStream,
		subscriptions: &sync.Map{},
		loggedIn:      false,
		username:      "<unknown>",
	}

	mh := wrapMessageHandler(mel, connInfo, messageHandler)

	connDead := make(chan bool)

	conn.SetCloseHandler(func(code int, text string) error {
		connDead <- true
		if connInfo.loggedIn {
			mel.RemoveConnection(connInfo.username, connInfo)
			log.WithFields(log.Fields{
				"addr":     conn.RemoteAddr().String(),
				"username": connInfo.username,
			}).Info("somebody has disconnected")
		}
		log.WithFields(log.Fields{
			"code": code,
			"text": text,
			"addr": conn.RemoteAddr().String(),
		}).Info("one of my connections is closed now")
		event := &MessageUserQuit{Username: connInfo.username}
		mel.IterateOverAllConnections(func(connInfo *ConnInfo) {
			connInfo.messageStream <- event
		})
		return nil
	})

	running := true

	// receiver
	go func() {
		for running {
			func() {
				defer func() {
					if err := recover(); err != nil {
						messageStream <- &MessageFatal{Message: fmt.Sprintf("%v", err)}
						log.WithFields(log.Fields{
							"addr": conn.RemoteAddr().String(),
							"name": connInfo.username,
							"err":  err,
						}).Error("panic while receiving a message")
						debug.PrintStack()
						running = false
					}
				}()
				if !running {
					return
				}
				var iface map[string]interface{}
				err := conn.ReadJSON(&iface)
				if err != nil {
					messageStream <- &MessageFatal{Message: "invalid JSON received"}
					log.WithFields(log.Fields{
						"addr": conn.RemoteAddr().String(),
						"name": connInfo.username,
						"err":  err,
					}).Error("cannot read a JSON message")
				}
				if !running {
					return
				}
				msg, err := LoadMessage(iface)
				if err != nil {
					messageStream <- &MessageFatal{Message: err.Error()}
					log.WithFields(log.Fields{
						"addr": conn.RemoteAddr().String(),
						"name": connInfo.username,
						"err":  err,
					}).Error("cannot process a JSON message")
					return
				}
				switch msg.(type) {
				case *MessageQuit:
					log.WithFields(log.Fields{
						"addr": conn.RemoteAddr().String(),
						"name": connInfo.username,
					}).Info("somebody wants to disconnect")
					running = false
					conn.Close()
				}
				go mh(msg)
			}()
		}
	}()

	// sender
	go func() {
		for running {
			select {
			case _ = <-connDead:
				running = false
			case msg := <-messageStream:
				func() {
					defer func() {
						if err := recover(); err != nil {
							log.WithFields(log.Fields{
								"addr": conn.RemoteAddr().String(),
								"name": connInfo.username,
								"err":  err,
							}).Error("panic while sending a message")
							debug.PrintStack()
							running = false
							conn.Close()
						}
					}()
					if !running {
						return
					}
					iface, err := MessageToIface(msg)
					if err != nil {
						log.WithFields(log.Fields{
							"addr": conn.RemoteAddr().String(),
							"name": connInfo.username,
							"err":  err,
						}).Error("cannot convert message to a map[string]interface{}")
						return
					}
					if !running {
						return
					}
					err = conn.WriteJSON(iface)
					if err != nil {
						log.WithFields(log.Fields{
							"addr": conn.RemoteAddr().String(),
							"name": connInfo.username,
							"err":  err,
						}).Error("unable to write JSON message")
					}
					switch msg.(type) {
					case *MessageQuit:
						log.WithFields(log.Fields{
							"addr": conn.RemoteAddr().String(),
							"name": connInfo.username,
						}).Info("disconnecting somebody")
						running = false
						conn.Close()
					case *MessageFatal:
						log.WithFields(log.Fields{
							"addr": conn.RemoteAddr().String(),
							"name": connInfo.username,
						}).Info("fatal error")
						running = false
						conn.Close()
					}
				}()
			}
		}
	}()
}
