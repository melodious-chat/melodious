package main

import "github.com/apex/log"

// Note: fatals are sent on database errors only in response to register and login

func handleRegisterMessage(mel *Melodious, connInfo *ConnInfo, message BaseMessage, send func(BaseMessage)) {
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
			send(&MessageNote{Message: "you are a server owner now"})
		}
	}
}

func handleLoginMessage(mel *Melodious, connInfo *ConnInfo, message BaseMessage, send func(BaseMessage)) {
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
}

func handleNewChannelMessage(mel *Melodious, connInfo *ConnInfo, message BaseMessage, send func(BaseMessage)) {
	cn, ct := message.(*MessageNewChannel).Name, message.(*MessageNewChannel).Topic
	nc := &MessageNewChannel{Name: cn, Topic: ct}
	owner, err := mel.Database.IsUserOwner(connInfo.username)
	if err != nil {
		send(&MessageFail{Message: "sorry, an internal database error has occured"})
		log.WithFields(log.Fields{
			"addr": connInfo.connection.RemoteAddr().String(),
			"name": connInfo.username,
			"err":  err,
		}).Error("error when checking if user is an owner")
	} else if owner {
		err = mel.Database.NewChannel(cn, ct)
		if err != nil {
			send(&MessageFail{Message: "sorry, an internal database error has occured"})
			log.WithFields(log.Fields{
				"addr": connInfo.connection.RemoteAddr().String(),
				"name": connInfo.username,
				"err":  err,
			}).Error("error when creating a channel")
		} else {
			send(&MessageOk{Message: "created a channel successfully"})
			mel.IterateOverAllConnections(func(connInfo *ConnInfo) {
				connInfo.messageStream <- nc
			})
		}
	} else {
		send(&MessageFail{Message: "no permissions"})
	}
}

func handleChannelTopicMessage(mel *Melodious, connInfo *ConnInfo, message BaseMessage, send func(BaseMessage)) {
	cn, ct := message.(*MessageChannelTopic).Name, message.(*MessageChannelTopic).Topic
	mct := &MessageChannelTopic{Name: cn, Topic: ct}
	owner, err := mel.Database.IsUserOwner(connInfo.username)
	if err != nil {
		send(&MessageFail{Message: "sorry, an internal database error has occured"})
		log.WithFields(log.Fields{
			"addr": connInfo.connection.RemoteAddr().String(),
			"name": connInfo.username,
			"err":  err,
		}).Error("error when checking if user is an owner")
	} else if owner {
		err = mel.Database.SetChannelTopic(cn, ct)
		if err != nil {
			send(&MessageFail{Message: "sorry, an internal database error has occured"})
			log.WithFields(log.Fields{
				"addr": connInfo.connection.RemoteAddr().String(),
				"name": connInfo.username,
				"err":  err,
			}).Error("error when changing channel topic")
		} else {
			send(&MessageOk{Message: "changed channel topic successfully"})
			mel.IterateOverAllConnections(func(connInfo *ConnInfo) {
				connInfo.messageStream <- mct
			})
		}
	} else {
		send(&MessageFail{Message: "no permissions"})
	}
}

func handleDeleteChannelMessage(mel *Melodious, connInfo *ConnInfo, message BaseMessage, send func(BaseMessage)) {
	cn := message.(*MessageDeleteChannel).Name
	dc := &MessageDeleteChannel{Name: cn}
	owner, err := mel.Database.IsUserOwner(connInfo.username)
	if err != nil {
		send(&MessageFail{Message: "sorry, an internal database error has occured"})
		log.WithFields(log.Fields{
			"addr": connInfo.connection.RemoteAddr().String(),
			"name": connInfo.username,
			"err":  err,
		}).Error("error when checking if user is an owner")
	} else if owner {
		err = mel.Database.DeleteChannel(cn)
		if err != nil {
			send(&MessageFail{Message: "sorry, an internal database error has occured"})
			log.WithFields(log.Fields{
				"addr": connInfo.connection.RemoteAddr().String(),
				"name": connInfo.username,
				"err":  err,
			}).Error("error when deleting a channel")
		} else {
			send(&MessageOk{Message: "deleted a channel successfully"})
			mel.IterateOverAllConnections(func(connInfo *ConnInfo) {
				connInfo.messageStream <- dc
			})
		}
	} else {
		send(&MessageFail{Message: "no permissions"})
	}
}

func handleQuitMessage(mel *Melodious, connInfo *ConnInfo, message BaseMessage, send func(BaseMessage)) {
	connInfo.connection.Close()
}

// messageHandler - handles messages received from users
func messageHandler(mel *Melodious, connInfo *ConnInfo, message BaseMessage, send func(BaseMessage)) {
	switch message.(type) {
	case *MessageRegister:
		handleRegisterMessage(mel, connInfo, message, send)
	case *MessageLogin:
		handleLoginMessage(mel, connInfo, message, send)
	case *MessageNewChannel:
		handleNewChannelMessage(mel, connInfo, message, send)
	case *MessageChannelTopic:
		handleChannelTopicMessage(mel, connInfo, message, send)
	case *MessageDeleteChannel:
		handleDeleteChannelMessage(mel, connInfo, message, send)
	case *MessageQuit:
		handleQuitMessage(mel, connInfo, message, send)
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
