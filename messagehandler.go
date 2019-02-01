package main

import "github.com/apex/log"

// messageHandler - handles messages received from users
func messageHandler(mel *Melodious, connInfo *ConnInfo, message BaseMessage, send func(BaseMessage)) {
	switch message.(type) {
	case *MessageRegister:
		if connInfo.loggedIn {
			send(&MessageFail{Message: "you are already logged in"})
			return
		}
		m := message.(*MessageRegister)
		exists, err := mel.Database.UserExists(m.Name)
		if err != nil {
			log.WithFields(log.Fields{
				"addr": connInfo.connection.RemoteAddr().String(),
				"name": m.Name,
				"err":  err,
			}).Error("error when checking if given user exists")
			send(&MessageFatal{Message: "sorry, an internal database error has occured"})
			return
		}
		if exists {
			send(&MessageFail{Message: "sorry, but there's already such a user with this nickname"})
			return
		}
		hasusers, err := mel.Database.HasUsers()
		firstrun := !hasusers
		if err != nil {
			log.WithFields(log.Fields{
				"addr": connInfo.connection.RemoteAddr().String(),
				"name": m.Name,
				"err":  err,
			}).Error("error when checking if database has users")
			send(&MessageFatal{Message: "sorry, an internal database error has occured"})
		} else if firstrun {
			err = mel.Database.RegisterUserOwner(m.Name, m.Pass, true)
		} else {
			err = mel.Database.RegisterUser(m.Name, m.Pass)
		}
		if err != nil {
			log.WithFields(log.Fields{
				"addr": connInfo.connection.RemoteAddr().String(),
				"name": m.Name,
				"err":  err,
			}).Error("error when registering a user")
			send(&MessageFatal{Message: "sorry, an internal database error has occured"})
		} else {
			connInfo.username = m.Name
			connInfo.loggedIn = true
			mel.PutConnection(m.Name, connInfo)
			log.WithFields(log.Fields{
				"addr": connInfo.connection.RemoteAddr().String(),
				"name": m.Name,
			}).Info("somebody has registered")
			if firstrun {
				log.Info("first run: registering as owner")
			}
			send(&MessageOk{Message: "done; you are now logged in"})
			if firstrun {
				send(&MessageNote{Message: "you are the server owner now"})
			}
		}
	case *MessageLogin:
		if connInfo.loggedIn {
			send(&MessageFail{Message: "you are already logged in"})
			return
		}
		m := message.(*MessageLogin)
		ok, err := mel.Database.CheckUserPassword(m.Name, m.Pass)
		if err != nil {
			send(&MessageFail{Message: err.Error()})
		} else if ok {
			connInfo.username = m.Name
			connInfo.loggedIn = true
			mel.PutConnection(m.Name, connInfo)
			log.WithFields(log.Fields{
				"addr": connInfo.connection.RemoteAddr().String(),
				"name": m.Name,
			}).Info("somebody has logged in")
			send(&MessageOk{Message: "done; you are now logged in"})
		}
	case *MessageQuit:
		connInfo.connection.Close()
	}
}

// wrapMessageHandler - wraps a message handler to allow passing it without explicitly passing some context-specific data
func wrapMessageHandler(
	mel *Melodious,
	connInfo *ConnInfo,
	f func(
		mel *Melodious,
		connInfo *ConnInfo,
		message BaseMessage,
		send func(BaseMessage),
	),
) func(BaseMessage) {
	return func(message BaseMessage) {
		send := func(m BaseMessage) {
			if id, ok := message.GetData().GetID(); ok {
				m.GetData().SetID(id)
			}
			connInfo.messageStream <- m
		}

		f(mel, connInfo, message, send)
	}
}
