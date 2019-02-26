package main

import (
	"database/sql"
	"fmt"
	"runtime/debug"
	"strconv"
	"strings"

	"github.com/apex/log"
)

// Note: fatals are sent on database errors only in response to register and login

func handleRegisterMessage(mel *Melodious, connInfo *ConnInfo, message BaseMessage, send func(BaseMessage)) {
	if connInfo.loggedIn {
		send(&MessageFail{Message: "you are already logged in"})
		return
	}
	m := message.(*MessageRegister)
	banned, err := mel.Database.IsUserBanned(m.Name, strings.Split(connInfo.connection.RemoteAddr().String(), ":")[0])
	if err != nil {
		send(&MessageFail{Message: err.Error()})
		return
	} else if banned {
		send(&MessageFatal{Message: "you are banned"})
		return
	}
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
		err = mel.Database.RegisterUserOwner(m.Name, m.Pass, true, strings.Split(connInfo.connection.RemoteAddr().String(), ":")[0])
	} else {
		err = mel.Database.RegisterUser(m.Name, m.Pass, strings.Split(connInfo.connection.RemoteAddr().String(), ":")[0])
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
	banned, err := mel.Database.IsUserBanned(m.Name, strings.Split(connInfo.connection.RemoteAddr().String(), ":")[0])
	if err != nil {
		send(&MessageFail{Message: err.Error()})
		return
	} else if banned {
		send(&MessageFatal{Message: "you are banned"})
		return
	}
	ok, err := mel.Database.CheckUserPassword(m.Name, m.Pass)
	if err != nil {
		send(&MessageFail{Message: err.Error()})
	} else if !ok {
		send(&MessageFatal{Message: "invalid credentials"})
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

	exists, err := mel.Database.ChannelExists(message.(*MessageSubscribe).Name)
	if err != nil {
		send(&MessageFail{Message: "sorry, an internal database error has occured"})
		log.WithFields(log.Fields{
			"addr": connInfo.connection.RemoteAddr().String(),
			"name": connInfo.username,
			"err":  err,
		}).Error("error when checking if the channel exists")
		return
	} else if !exists {
		send(&MessageFail{Message: "no such channel"})
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
	can, err := connInfo.HasPerm(message.(*MessagePostMsg).Channel, "perms.post-message")
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

	if _, ok := connInfo.subscriptions.Load(message.(*MessagePostMsg).Channel); !ok {
		send(&MessageFail{Message: "not subscribed to the sending channel"})
		return
	}
	author := connInfo.username
	ids := scanForPings(message.(*MessagePostMsg).Content)
	pings := []string{}
	unknownids := []int{}
	for _, id := range ids {
		user, err := mel.Database.GetUser(id)
		if err == sql.ErrNoRows {
			unknownids = append(unknownids, id)
		} else if err != nil {
			send(&MessageFail{Message: "sorry, an internal database error has occured"})
			log.WithFields(log.Fields{
				"addr": connInfo.connection.RemoteAddr().String(),
				"name": connInfo.username,
				"err":  err,
			}).Error("error when getting user info")
		} else {
			pings = append(pings, user.Username)
		}
	}
	msg, err := mel.Database.PostMessage(message.(*MessagePostMsg).Channel, message.(*MessagePostMsg).Content, pings, author)
	if err != nil {
		send(&MessageFail{Message: "sorry, an internal database error has occured"})
		log.WithFields(log.Fields{
			"addr": connInfo.connection.RemoteAddr().String(),
			"name": connInfo.username,
			"err":  err,
		}).Error("error when posting a message")
	} else {
		im := &MessagePostMsg{Channel: message.(*MessagePostMsg).Channel, MsgObj: msg}
		if len(pings) != 0 {
			ping := &MessagePing{Message: im.MsgObj, Channel: im.Channel}
			for _, username := range pings {
				mel.IterateOverConnections(username, func(connInfo *ConnInfo) {
					connInfo.messageStream <- ping
				})
			}
		}
		mel.IterateOverAllConnections(func(connInfo *ConnInfo) {
			if subbed, ok := connInfo.subscriptions.Load(message.(*MessagePostMsg).Channel); subbed == true && ok {
				connInfo.messageStream <- im
			}
		})
		unkidstr := ""
		for _, id := range unknownids {
			unkidstr += strconv.Itoa(id) + " "
		}
		if len(unknownids) != 0 {
			send(&MessageNote{Message: "warning: unknown ids " + unkidstr})
		}
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
	} else if !can {
		send(&MessageFail{Message: "no permissions"})
		return
	}

	username := message.(*MessageKick).Username
	if username == connInfo.username {
		send(&MessageFail{Message: "you can't kick or ban yourself"})
		return
	}

	exists, err := mel.Database.UserExists(username)
	if err != nil {
		send(&MessageFail{Message: "sorry, an internal database error has occured"})
		log.WithFields(log.Fields{
			"addr": connInfo.connection.RemoteAddr().String(),
			"name": connInfo.username,
			"err":  err,
		}).Error("error when checking if a user exists")
		return
	} else if !exists {
		send(&MessageFail{Message: "no such user"})
		return
	}
	mel.IterateOverConnections(username, func(connInfo *ConnInfo) {
		connInfo.messageStream <- &MessageFatal{Message: "you've been kicked or banned"}
	})
	if message.(*MessageKick).Ban {
		err = mel.Database.Ban(username)
		if err != nil {
			send(&MessageFail{Message: "sorry, an internal database error has occured"})
			log.WithFields(log.Fields{
				"addr": connInfo.connection.RemoteAddr().String(),
				"name": connInfo.username,
				"err":  err,
			}).Error("error when banning a user")
			return
		}
		send(&MessageOk{Message: "kicked and banned user " + username})
		return
	}
	send(&MessageOk{Message: "kicked user " + username})
}

func handleNewGroupMessage(mel *Melodious, connInfo *ConnInfo, message BaseMessage, send func(BaseMessage)) {
	can, err := mel.Database.IsUserOwner(connInfo.username)
	if err != nil {
		send(&MessageFail{Message: "sorry, an internal database error has occured"})
		log.WithFields(log.Fields{
			"addr": connInfo.connection.RemoteAddr().String(),
			"name": connInfo.username,
			"err":  err,
		}).Error("error when checking if user can manage groups")
		return
	} else if !can {
		send(&MessageFail{Message: "no permissions"})
		return
	}
	_, err = mel.Database.AddGroup(message.(*MessageNewGroup).Name)
	if err != nil {
		send(&MessageFail{Message: "sorry, an internal database error has occured"})
		log.WithFields(log.Fields{
			"addr": connInfo.connection.RemoteAddr().String(),
			"name": connInfo.username,
			"err":  err,
		}).Error("error when adding a group")
		return
	}
	send(&MessageOk{Message: "created group " + message.(*MessageNewGroup).Name})
}

func handleDeleteGroupMessage(mel *Melodious, connInfo *ConnInfo, message BaseMessage, send func(BaseMessage)) {
	can, err := mel.Database.IsUserOwner(connInfo.username)
	if err != nil {
		send(&MessageFail{Message: "sorry, an internal database error has occured"})
		log.WithFields(log.Fields{
			"addr": connInfo.connection.RemoteAddr().String(),
			"name": connInfo.username,
			"err":  err,
		}).Error("error when checking if user can manage groups")
		return
	} else if !can {
		send(&MessageFail{Message: "no permissions"})
		return
	}
	exists, err := mel.Database.GroupExists(message.(*MessageDeleteGroup).Name)
	if err != nil {
		send(&MessageFail{Message: "sorry, an internal database error has occured"})
		log.WithFields(log.Fields{
			"addr": connInfo.connection.RemoteAddr().String(),
			"name": connInfo.username,
			"err":  err,
		}).Error("error when checking if a group exists")
		return
	} else if !exists {
		send(&MessageFail{Message: "no such group"})
		return
	}
	err = mel.Database.DeleteGroup(message.(*MessageDeleteGroup).Name)
	if err != nil {
		send(&MessageFail{Message: "sorry, an internal database error has occured"})
		log.WithFields(log.Fields{
			"addr": connInfo.connection.RemoteAddr().String(),
			"name": connInfo.username,
			"err":  err,
		}).Error("error when deleting a group")
		return
	}
	send(&MessageOk{Message: "deleted group " + message.(*MessageDeleteGroup).Name})
}

func handleSetFlagMessage(mel *Melodious, connInfo *ConnInfo, message BaseMessage, send func(BaseMessage)) {
	can, err := mel.Database.IsUserOwner(connInfo.username)
	if err != nil {
		send(&MessageFail{Message: "sorry, an internal database error has occured"})
		log.WithFields(log.Fields{
			"addr": connInfo.connection.RemoteAddr().String(),
			"name": connInfo.username,
			"err":  err,
		}).Error("error when checking if user is owner")
		return
	} else if !can {
		send(&MessageFail{Message: "no permissions"})
		return
	}
	exists, err := mel.Database.GroupExists(message.(*MessageSetFlag).Group)
	if err != nil {
		send(&MessageFail{Message: "sorry, an internal database error has occured"})
		log.WithFields(log.Fields{
			"addr": connInfo.connection.RemoteAddr().String(),
			"name": connInfo.username,
			"err":  err,
		}).Error("error when checking if a group exists")
		return
	} else if !exists {
		send(&MessageFail{Message: "no such group"})
		return
	}
	procmsg := message.(*MessageSetFlag)
	_, err = mel.Database.SetFlag(&Flag{HasID: false, Group: procmsg.Group, Name: procmsg.Name, Flag: procmsg.Flag})
	if err != nil {
		send(&MessageFail{Message: "sorry, an internal database error has occured"})
		log.WithFields(log.Fields{
			"addr": connInfo.connection.RemoteAddr().String(),
			"name": connInfo.username,
			"err":  err,
		}).Error("error when setting a flag")
		return
	}
	send(&MessageOk{Message: "set flag " + procmsg.Name + " for group " + procmsg.Group})
}

func handleDeleteFlagMessage(mel *Melodious, connInfo *ConnInfo, message BaseMessage, send func(BaseMessage)) {
	can, err := mel.Database.IsUserOwner(connInfo.username)
	if err != nil {
		send(&MessageFail{Message: "sorry, an internal database error has occured"})
		log.WithFields(log.Fields{
			"addr": connInfo.connection.RemoteAddr().String(),
			"name": connInfo.username,
			"err":  err,
		}).Error("error when checking if user is owner")
		return
	} else if !can {
		send(&MessageFail{Message: "no permissions"})
		return
	}
	exists, err := mel.Database.GroupExists(message.(*MessageDeleteFlag).Group)
	if err != nil {
		send(&MessageFail{Message: "sorry, an internal database error has occured"})
		log.WithFields(log.Fields{
			"addr": connInfo.connection.RemoteAddr().String(),
			"name": connInfo.username,
			"err":  err,
		}).Error("error when checking if a group exists")
		return
	} else if !exists {
		send(&MessageFail{Message: "no such group"})
		return
	}
	procmsg := message.(*MessageDeleteFlag)
	err = mel.Database.DeleteFlag(&Flag{HasID: false, Group: procmsg.Group, Name: procmsg.Name})
	if err != nil {
		send(&MessageFail{Message: "sorry, an internal database error has occured"})
		log.WithFields(log.Fields{
			"addr": connInfo.connection.RemoteAddr().String(),
			"name": connInfo.username,
			"err":  err,
		}).Error("error when removing a flag")
		return
	}
	send(&MessageOk{Message: "removed flag " + procmsg.Name + " from group " + procmsg.Group})
}

func handleTypingMessage(mel *Melodious, connInfo *ConnInfo, message BaseMessage, send func(BaseMessage)) {
	can, err := connInfo.HasPerm(message.(*MessageTyping).Channel, "perms.post-message")
	if err != nil {
		send(&MessageFail{Message: "sorry, an internal database error has occured"})
		log.WithFields(log.Fields{
			"addr": connInfo.connection.RemoteAddr().String(),
			"name": connInfo.username,
			"err":  err,
		}).Error("error when checking if user can type")
		return
	}
	if !can {
		send(&MessageFail{Message: "no permissions"})
		return
	}

	if message.(*MessageTyping).HasUsername {
		send(&MessageNote{Message: "you cannot set username field in typing message"})
	}
	username := connInfo.username
	procmsg := message.(*MessageTyping)
	mt := &MessageTyping{Channel: procmsg.Channel, Username: username, Typing: procmsg.Typing}
	mel.IterateOverAllConnections(func(connInfo *ConnInfo) {
		if subbed, ok := connInfo.subscriptions.Load(message.(*MessageTyping).Channel); subbed == true && ok {
			connInfo.messageStream <- mt
		}
	})
	if procmsg.Typing {
		send(&MessageOk{Message: "typing in " + procmsg.Channel})
	} else {
		send(&MessageOk{Message: "stopped typing in " + procmsg.Channel})
	}
}

func handleNewGroupHolderMessage(mel *Melodious, connInfo *ConnInfo, message BaseMessage, send func(BaseMessage)) {
	can, err := mel.Database.IsUserOwner(connInfo.username)
	if err != nil {
		send(&MessageFail{Message: "sorry, an internal database error has occured"})
		log.WithFields(log.Fields{
			"addr": connInfo.connection.RemoteAddr().String(),
			"name": connInfo.username,
			"err":  err,
		}).Error("error when checking if user is owner")
		return
	} else if !can {
		send(&MessageFail{Message: "no permissions"})
		return
	}

	procmsg := message.(*MessageNewGroupHolder)
	if procmsg.Channel != "" {
		exists, err := mel.Database.ChannelExists(procmsg.Channel)
		if err != nil {
			send(&MessageFail{Message: "sorry, an internal database error has occured"})
			log.WithFields(log.Fields{
				"addr": connInfo.connection.RemoteAddr().String(),
				"name": connInfo.username,
				"err":  err,
			}).Error("error when checking if the channel exists")
			return
		} else if !exists {
			send(&MessageFail{Message: "no such channel"})
			return
		}
	}
	gh := &GroupHolder{Group: procmsg.Group, User: procmsg.User, Channel: procmsg.Channel}
	_, err = mel.Database.AddGroupHolder(gh)
	if err != nil {
		send(&MessageFail{Message: "sorry, an internal database error has occured"})
		log.WithFields(log.Fields{
			"addr": connInfo.connection.RemoteAddr().String(),
			"name": connInfo.username,
			"err":  err,
		}).Error("error when adding a group holder")
		return
	}
	sendmsg := &MessageOk{Message: ""}
	if procmsg.User != "" {
		sendmsg.Message += "assigned user " + procmsg.User
	} else {
		sendmsg.Message += "assigned everyone"
	}
	sendmsg.Message += " to group " + procmsg.Group
	if procmsg.Channel != "" {
		sendmsg.Message += " to channel " + procmsg.Channel
	}
	send(sendmsg)
}

func handleDeleteGroupHolderMessage(mel *Melodious, connInfo *ConnInfo, message BaseMessage, send func(BaseMessage)) {
	can, err := mel.Database.IsUserOwner(connInfo.username)
	if err != nil {
		send(&MessageFail{Message: "sorry, an internal database error has occured"})
		log.WithFields(log.Fields{
			"addr": connInfo.connection.RemoteAddr().String(),
			"name": connInfo.username,
			"err":  err,
		}).Error("error when checking if user is owner")
		return
	} else if !can {
		send(&MessageFail{Message: "no permissions"})
		return
	}

	procmsg := message.(*MessageDeleteGroupHolder)
	err = mel.Database.DeleteGroupHolder(procmsg.ID)
	if err != nil {
		send(&MessageFail{Message: "sorry, an internal database error has occured"})
		log.WithFields(log.Fields{
			"addr": connInfo.connection.RemoteAddr().String(),
			"name": connInfo.username,
			"err":  err,
		}).Error("error when deleting a group holder")
		return
	}
	send(&MessageOk{Message: "deleted group holder with id " + string(procmsg.ID)})
}

func handleDeleteMsgMessage(mel *Melodious, connInfo *ConnInfo, message BaseMessage, send func(BaseMessage)) {
	procmsg := message.(*MessageDeleteMsg)
	channel, msg, err := mel.Database.GetMessageDetails(procmsg.ID)
	if err == sql.ErrNoRows {
		send(&MessageFail{Message: "no such message with id " + strconv.Itoa(procmsg.ID)})
		return
	} else if err != nil {
		send(&MessageFail{Message: "sorry, an internal database error has occured"})
		log.WithFields(log.Fields{
			"addr": connInfo.connection.RemoteAddr().String(),
			"name": connInfo.username,
			"err":  err,
		}).Error("error when fetching message details")
		return
	}
	if msg.Author != connInfo.username {
		can, err := connInfo.HasPerm(channel, "perms.delete-message")
		if err != nil {
			send(&MessageFail{Message: "sorry, an internal database error has occured"})
			log.WithFields(log.Fields{
				"addr": connInfo.connection.RemoteAddr().String(),
				"name": connInfo.username,
				"err":  err,
			}).Error("error when checking if user has permissions")
			return
		} else if !can {
			send(&MessageFail{Message: "no permissions"})
			return
		}
	}
	err = mel.Database.DeleteMessage(procmsg.ID)
	if err != nil {
		send(&MessageFail{Message: "sorry, an internal database error has occured"})
		log.WithFields(log.Fields{
			"addr": connInfo.connection.RemoteAddr().String(),
			"name": connInfo.username,
			"err":  err,
		}).Error("error when deleting a message")
		return
	}
	send(&MessageOk{Message: "deleted message with id " + strconv.Itoa(procmsg.ID)})
}

func handleGetGroupHoldersMessage(mel *Melodious, connInfo *ConnInfo, message BaseMessage, send func(BaseMessage)) {
	can, err := mel.Database.IsUserOwner(connInfo.username)
	if err != nil {
		send(&MessageFail{Message: "sorry, an internal database error has occured"})
		log.WithFields(log.Fields{
			"addr": connInfo.connection.RemoteAddr().String(),
			"name": connInfo.username,
			"err":  err,
		}).Error("error when checking if user is owner")
		return
	} else if !can {
		send(&MessageFail{Message: "no permissions"})
		return
	}
	ghs, err := mel.Database.GetGroupHolders()
	if err != nil {
		send(&MessageFail{Message: "sorry, an internal database error has occured"})
		log.WithFields(log.Fields{
			"addr": connInfo.connection.RemoteAddr().String(),
			"name": connInfo.username,
			"err":  err,
		}).Error("error when getting group holders")
		return
	}
	send(&MessageGetGroupHolders{GroupHolders: ghs})
}

// messageHandler - handles messages received from users
func messageHandler(mel *Melodious, connInfo *ConnInfo, message BaseMessage, send func(BaseMessage)) {
	defer func() {
		if err := recover(); err != nil {
			send(&MessageFatal{Message: fmt.Sprintf("%v", err)})
			log.WithFields(log.Fields{
				"addr": connInfo.connection.RemoteAddr().String(),
				"name": connInfo.username,
				"err":  err,
			}).Error("panic while receiving a message")
			debug.PrintStack()
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
		case *MessageNewGroup:
			handleNewGroupMessage(mel, connInfo, message, send)
		case *MessageDeleteGroup:
			handleDeleteGroupMessage(mel, connInfo, message, send)
		case *MessageSetFlag:
			handleSetFlagMessage(mel, connInfo, message, send)
		case *MessageDeleteFlag:
			handleDeleteFlagMessage(mel, connInfo, message, send)
		case *MessageTyping:
			handleTypingMessage(mel, connInfo, message, send)
		case *MessageNewGroupHolder:
			handleNewGroupHolderMessage(mel, connInfo, message, send)
		case *MessageDeleteMsg:
			handleDeleteMsgMessage(mel, connInfo, message, send)
		case *MessageGetGroupHolders:
			handleGetGroupHoldersMessage(mel, connInfo, message, send)
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
