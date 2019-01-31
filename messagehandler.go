package main

import "github.com/apex/log"

// messageHandler - handles messages received from users
func messageHandler(mel *Melodious, connInfo *ConnInfo, message BaseMessage) {
	switch message.(type) {
	case *MessageRegister:
		if connInfo.loggedIn {
			connInfo.messageStream <- &MessageFail{Message: "you are already logged in"}
			return
		}
		m := message.(*MessageRegister)
		err := mel.Database.RegisterUser(m.Name, m.Pass)
		if err != nil {
			connInfo.messageStream <- &MessageFail{Message: err.Error()}
		} else {
			connInfo.username = m.Name
			connInfo.loggedIn = true
			mel.PutConnection(m.Name, connInfo)
			log.WithFields(log.Fields{"username": m.Name}).Info("somebody has registered")
			connInfo.messageStream <- &MessageOk{Message: "done; you are now logged in"}
		}
	case *MessageLogin:
		if connInfo.loggedIn {
			connInfo.messageStream <- &MessageFail{Message: "you are already logged in"}
			return
		}
		m := message.(*MessageLogin)
		ok, err := mel.Database.CheckUserPassword(m.Name, m.Pass)
		if err != nil {
			connInfo.messageStream <- &MessageFail{Message: err.Error()}
		} else if ok {
			connInfo.username = m.Name
			connInfo.loggedIn = true
			mel.PutConnection(m.Name, connInfo)
			log.WithFields(log.Fields{"username": m.Name}).Info("somebody has logged in")
			connInfo.messageStream <- &MessageOk{Message: "done; you are now logged in"}
		}
	case *MessageQuit:
		if connInfo.loggedIn {
			connInfo.loggedIn = false
		}
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
	),
) func(BaseMessage) {
	return func(message BaseMessage) {
		f(mel, connInfo, message)
	}
}
