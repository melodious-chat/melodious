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
		fmt.Printf(msg.GetType() + " %U\n", msg)
		handler(mel, messageStream, msg)
	}

	connDead := make(chan bool)

	conn.SetCloseHandler(func(code int, text string) error {
		connDead <- true
		log.WithFields(log.Fields{"code": code, "text": text}).Info("one of my connections is closed now")
		return nil
	})

	// receiver
	go func() {
		for {
			var iface map[string]interface{}
			err := conn.ReadJSON(&iface)
			if err != nil {
				log.WithField("err", err).Error("cannot read a JSON message")
				continue
			}
			msg, err := LoadMessage(iface)
			if err != nil {
				log.WithField("err", err).Error("cannot process a JSON message")
				continue
			}
			go messageHandler(msg)
		}
	}()

	// sender
	go func() {
		for {
			select {
			case _ = <-connDead:
				break
			case msg := <-messageStream:
				iface, err := MessageToIface(msg)
				if err != nil {
					log.WithField("err", err).Error("cannot convert message to a map[string]interface{}")
					continue
				}
				err = conn.WriteJSON(iface)
				if err != nil {
					log.WithField("err", err).Error("unable to write JSON message")
				}
			}
		}
	}()
}
