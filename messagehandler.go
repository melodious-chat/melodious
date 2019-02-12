package main

import (
	"fmt"

	"github.com/apex/log"
)

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
		event := &MessageRegister{Name: m.Name}
		mel.IterateOverAllConnections(func(connInfo *ConnInfo) {
			connInfo.messageStream <- event
		})
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
		event := &MessageLogin{Name: m.Name}
		mel.IterateOverAllConnections(func(connInfo *ConnInfo) {
			connInfo.messageStream <- event
		})
	}
}

func handleNewChannelMessage(mel *Melodious, connInfo *ConnInfo, message BaseMessage, send func(BaseMessage)) {
	cn, ct := message.(*MessageNewChannel).Name, message.(*MessageNewChannel).Topic
	nc := &MessageNewChannel{Name: cn, Topic: ct}
	can, err := connInfo.HasPerm(cn, "perms.new-channel")
	if err != nil {
		send(&MessageFail{Message: "sorry, an internal database error has occured"})
		log.WithFields(log.Fields{
			"addr": connInfo.connection.RemoteAddr().String(),
			"name": connInfo.username,
			"err":  err,
		}).Error("error when checking if user can create a channel")
	} else if can {
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
	can, err := connInfo.HasPerm(cn, "perms.channel-topic")
	if err != nil {
		send(&MessageFail{Message: "sorry, an internal database error has occured"})
		log.WithFields(log.Fields{
			"addr": connInfo.connection.RemoteAddr().String(),
			"name": connInfo.username,
			"err":  err,
		}).Error("error when checking if user can change channel topic")
	} else if can {
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
	can, err := connInfo.HasPerm(cn, "perms.delete-channel")
	if err != nil {
		send(&MessageFail{Message: "sorry, an internal database error has occured"})
		log.WithFields(log.Fields{
			"addr": connInfo.connection.RemoteAddr().String(),
			"name": connInfo.username,
			"err":  err,
		}).Error("error when checking if user can delete a channel")
	} else if can {
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
	event := &MessageUserQuit{Username: connInfo.username}
	mel.IterateOverAllConnections(func(connInfo *ConnInfo) {
		connInfo.messageStream <- event
	})
}

func handleSubscribeMessage(mel *Melodious, connInfo *ConnInfo, message BaseMessage, send func(BaseMessage)) {
	can, err := connInfo.HasPerm(message.(*MessageSubscribe).Name, "perms.subscribe")
	if err != nil {
		send(&MessageFail{Message: "sorry, an internal database error has occured"})
		log.WithFields(log.Fields{
			"addr": connInfo.connection.RemoteAddr().String(),
			"name": connInfo.username,
			"err":  err,
		}).Error("error when checking if user can (un)subscribe to a channel")
		return
	}
	if !can {
		send(&MessageFail{Message: "no permissions"})
		return
	}

	if message.(*MessageSubscribe).Subbed {
		connInfo.subscriptions.Store(message.(*MessageSubscribe).Name, message.(*MessageSubscribe).Subbed)
		send(&MessageOk{Message: "subscribed to channel " + message.(*MessageSubscribe).Name})
	} else {
		connInfo.subscriptions.Delete(message.(*MessageSubscribe).Name)
		send(&MessageOk{Message: "unsubscribed from channel " + message.(*MessageSubscribe).Name})
	}
}

func handlePostMsgMessage(mel *Melodious, connInfo *ConnInfo, message BaseMessage, send func(BaseMessage)) {
	can, err := connInfo.HasPerm(message.(*MessageSubscribe).Name, "perms.post-message")
	if err != nil {
		send(&MessageFail{Message: "sorry, an internal database error has occured"})
		log.WithFields(log.Fields{
			"addr": connInfo.connection.RemoteAddr().String(),
			"name": connInfo.username,
			"err":  err,
		}).Error("error when checking if user can post a message")
		return
	}
	if !can {
		send(&MessageFail{Message: "no permissions"})
		return
	}

	if message.(*MessagePostMsg).HasAuthor {
		send(&MessageNote{Message: "you cannot set author field in post-message message"})
	}
	author := connInfo.username
	err = mel.Database.PostMessage(message.(*MessagePostMsg).Channel, message.(*MessagePostMsg).Content, []string{""}, author) // todo pings
	if err != nil {
		send(&MessageFail{Message: "sorry, an internal database error has occured"})
		log.WithFields(log.Fields{
			"addr": connInfo.connection.RemoteAddr().String(),
			"name": connInfo.username,
			"err":  err,
		}).Error("error when posting a message")
	} else {
		im := &MessagePostMsg{Content: message.(*MessagePostMsg).Content, Channel: message.(*MessagePostMsg).Channel, HasAuthor: true, Author: author}
		mel.IterateOverAllConnections(func(connInfo *ConnInfo) {
			if subbed, ok := connInfo.subscriptions.Load(message.(*MessagePostMsg).Channel); subbed == true && ok {
				connInfo.messageStream <- im
			}
		})
		send(&MessageOk{Message: "message sent to channel " + message.(*MessagePostMsg).Channel})
	}
}

func handleGetMsgsMessage(mel *Melodious, connInfo *ConnInfo, message BaseMessage, send func(BaseMessage)) {
	request := message.(*MessageGetMsgs)

	can, err := connInfo.HasPermChID(request.ChannelID, "perms.get-messages")
	if err != nil {
		send(&MessageFail{Message: "sorry, an internal database error has occured"})
		log.WithFields(log.Fields{
			"addr": connInfo.connection.RemoteAddr().String(),
			"name": connInfo.username,
			"err":  err,
		}).Error("error when checking if user can get messages")
		return
	}
	if !can {
		send(&MessageFail{Message: "no permissions"})
		return
	}

	msgs, err := mel.Database.GetMessages(request.ChannelID, request.MessageID, request.Amount)
	if err != nil {
		send(&MessageFail{Message: "sorry, an internal database error has occured"})
		log.WithFields(log.Fields{
			"addr": connInfo.connection.RemoteAddr().String(),
			"name": connInfo.username,
			"err":  err,
		}).Error("error when fetching messages")
	} else {
		send(&MessageGetMsgsResult{Messages: msgs})
	}
}

func handleListChannelsMessage(mel *Melodious, connInfo *ConnInfo, message BaseMessage, send func(BaseMessage)) {
	can, err := connInfo.HasPerm("", "perms.list-channels")
	if err != nil {
		send(&MessageFail{Message: "sorry, an internal database error has occured"})
		log.WithFields(log.Fields{
			"addr": connInfo.connection.RemoteAddr().String(),
			"name": connInfo.username,
			"err":  err,
		}).Error("error when checking if user can list channels")
		return
	}
	if !can {
		send(&MessageFail{Message: "no permissions"})
		return
	}

	if message.(*MessageListChannels).HasChannels {
		send(&MessageNote{Message: "you cannot set channels field in list-channels message"})
	}
	channels, err := mel.Database.ListChannels()
	if err != nil {
		send(&MessageFail{Message: "sorry, an internal database error has occured"})
		log.WithFields(log.Fields{
			"addr": connInfo.connection.RemoteAddr().String(),
			"name": connInfo.username,
			"err":  err,
		}).Error("error when listing channels")
	} else {
		send(&MessageListChannels{Channels: channels, HasChannels: true})
	}
}

func handleListUsersMessage(mel *Melodious, connInfo *ConnInfo, message BaseMessage, send func(BaseMessage)) {
	can, err := connInfo.HasPerm("", "perms.list-users")
	if err != nil {
		send(&MessageFail{Message: "sorry, an internal database error has occured"})
		log.WithFields(log.Fields{
			"addr": connInfo.connection.RemoteAddr().String(),
			"name": connInfo.username,
			"err":  err,
		}).Error("error when checking if user can list users")
		return
	}
	if !can {
		send(&MessageFail{Message: "no permissions"})
		return
	}

	if message.(*MessageListUsers).HasUsers {
		send(&MessageNote{Message: "you cannot set users field in list-users message"})
	}
	users, err := mel.Database.GetUsersList()
	if err != nil {
		send(&MessageFail{Message: "sorry, an internal database error has occured"})
		log.WithFields(log.Fields{
			"addr": connInfo.connection.RemoteAddr().String(),
			"name": connInfo.username,
			"err":  err,
		}).Error("error when getting users list")
	} else {
		statuses := []*UserStatus{}
		for _, user := range users {
			_, online := mel.UserConns.Load(user.Username)
			statuses = append(statuses, &UserStatus{User: user, Online: online})
		}
		send(&MessageListUsers{Users: statuses, HasUsers: true})
	}
}

func handleKickMessage(mel *Melodious, connInfo *ConnInfo, message BaseMessage, send func(BaseMessage)) {
	can, err := connInfo.HasPerm("", "perms.kickban")
	if err != nil {
		send(&MessageFail{Message: "sorry, an internal database error has occured"})
		log.WithFields(log.Fields{
			"addr": connInfo.connection.RemoteAddr().String(),
			"name": connInfo.username,
			"err":  err,
		}).Error("error when checking if user can kick and ban")
		return
	}
	if !can {
		send(&MessageFail{Message: "no permissions"})
		return
	}

	username := message.(*MessageKick).Username

	mel.IterateOverConnections(username, func(connInfo *ConnInfo) {
		connInfo.messageStream <- &MessageFatal{Message: "you've been kicked or banned"}
	})
	if message.(*MessageKick).Ban {
		err = mel.Database.DeleteUser(username)
		if err != nil {
			log.WithFields(log.Fields{
				"addr": connInfo.connection.RemoteAddr().String(),
				"name": connInfo.username,
				"err":  err,
			}).Error("error when banning a user")
			return
		}
		err = mel.Database.InsertBan(username, connInfo.connection.RemoteAddr().String())
		if err != nil {
			log.WithFields(log.Fields{
				"addr": connInfo.connection.RemoteAddr().String(),
				"name": connInfo.username,
				"err":  err,
			}).Error("error when banning a user")
			return
		}
	}
}

// messageHandler - handles messages received from users
func messageHandler(mel *Melodious, connInfo *ConnInfo, message BaseMessage, send func(BaseMessage)) {
	defer func() {
		if err := recover(); err != nil {
			send(&MessageFail{Message: fmt.Sprintf("%U", err)})
			log.WithFields(log.Fields{
				"addr": connInfo.connection.RemoteAddr().String(),
				"name": connInfo.username,
				"err":  err,
			}).Error("panic while receiving a message")
			connInfo.connection.Close()
		}
	}()

	if !connInfo.loggedIn {
		switch message.(type) {
		case *MessageRegister:
			handleRegisterMessage(mel, connInfo, message, send)
		case *MessageLogin:
			handleLoginMessage(mel, connInfo, message, send)
		}
	} else {
		switch message.(type) {
		case *MessageNewChannel:
			handleNewChannelMessage(mel, connInfo, message, send)
		case *MessageChannelTopic:
			handleChannelTopicMessage(mel, connInfo, message, send)
		case *MessageDeleteChannel:
			handleDeleteChannelMessage(mel, connInfo, message, send)
		case *MessageQuit:
			handleQuitMessage(mel, connInfo, message, send)
		case *MessageSubscribe:
			handleSubscribeMessage(mel, connInfo, message, send)
		case *MessagePostMsg:
			handlePostMsgMessage(mel, connInfo, message, send)
		case *MessageGetMsgs:
			handleGetMsgsMessage(mel, connInfo, message, send)
		case *MessageListChannels:
			handleListChannelsMessage(mel, connInfo, message, send)
		case *MessageListUsers:
			handleListUsersMessage(mel, connInfo, message, send)
		case *MessageKick:
			handleKickMessage(mel, connInfo, message, send)
		}
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
