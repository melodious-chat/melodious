package main

import (
	"fmt"

	"github.com/apex/log"
	"github.com/gorilla/websocket"
)

// handleConnection Handles users which are connected to Melodious
func handleConnection(mel *Melodious, conn *websocket.Conn) {

	messageStream := make(chan BaseMessage)

	messageHandler := func(msg BaseMessage) {
		fmt.Printf(msg.GetType()+" %U\n", msg)
		messageHandler(mel, messageStream, msg)
	}

	connDead := make(chan bool)

	conn.SetCloseHandler(func(code int, text string) error {
		connDead <- true
		log.WithFields(log.Fields{"code": code, "text": text}).Info("one of my connections is closed now")
		return nil
	})

	running := true

	// receiver
	go func() {
		for running {
			func() {
				defer func() {
					if err := recover(); err != nil {
						messageStream <- &MessageFail{Message: fmt.Sprintf("%U", err)}
						log.WithField("err", err).Error("panic while receiving a message")
						running = false
					}
				}()
				if !running {
					return
				}
				var iface map[string]interface{}
				err := conn.ReadJSON(&iface)
				if err != nil {
					log.WithField("err", err).Error("cannot read a JSON message")
				}
				if !running {
					return
				}
				msg, err := LoadMessage(iface)
				if err != nil {
					messageStream <- &MessageFail{Message: err.Error()}
					log.WithField("err", err).Error("cannot process a JSON message")
					return
				}
				go messageHandler(msg)
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
							log.WithField("err", err).Error("panic while sending a message")
							running = false
						}
					}()
					if !running {
						return
					}
					iface, err := MessageToIface(msg)
					if err != nil {
						log.WithField("err", err).Error("cannot convert message to a map[string]interface{}")
						return
					}
					if !running {
						return
					}
					err = conn.WriteJSON(iface)
					if err != nil {
						log.WithField("err", err).Error("unable to write JSON message")
					}
				}()
			}
		}
	}()
}
