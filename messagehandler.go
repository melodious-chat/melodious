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
		exists, err := mel.Database.UserExists(m.Name)
		if err != nil {
			log.WithFields(log.Fields{
				"addr": connInfo.connection.RemoteAddr().String(),
				"name": m.Name,
				"err":  err,
			}).Error("error when checking if given user exists")
			connInfo.messageStream <- &MessageFatal{Message: "sorry, an internal database error has occured"}
			return
		}
		if exists {
			connInfo.messageStream <- &MessageFail{Message: "sorry, but there's already such a user with this nickname"}
			return
		}
		hasusers, err := mel.Database.HasUsers()
		firstrun := !hasusers
		log.WithFields(log.Fields{"firstrun":firstrun,"hasusers":hasusers,"err":err}).Info("")
		if err != nil {
			log.WithFields(log.Fields{
				"addr": connInfo.connection.RemoteAddr().String(),
				"name": m.Name,
				"err":  err,
			}).Error("error when checking if database has users")
			connInfo.messageStream <- &MessageFatal{Message: "sorry, an internal database error has occured"}
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
			connInfo.messageStream <- &MessageFatal{Message: "sorry, an internal database error has occured"}
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
			connInfo.messageStream <- &MessageOk{Message: "done; you are now logged in"}
			if firstrun {
				connInfo.messageStream <- &MessageNote{Message: "you became the server owner"}
			}
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
			log.WithFields(log.Fields{
				"addr": connInfo.connection.RemoteAddr().String(),
				"name": m.Name,
			}).Info("somebody has logged in")
			connInfo.messageStream <- &MessageOk{Message: "done; you are now logged in"}
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
	),
) func(BaseMessage) {
	return func(message BaseMessage) {
		f(mel, connInfo, message)
	}
}
