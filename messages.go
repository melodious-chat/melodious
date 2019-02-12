package main

import (
	"errors"
)

// BaseMessage - A base struct for all messages
type BaseMessage interface {
	GetType() string
	GetData() *MessageData
}

// MessageData - Some extra message data
type MessageData struct {
	id    string
	hasID bool
}

// SetID - sets message ID
func (m *MessageData) SetID(id string) {
	m.hasID = true
	m.id = id
}

// GetID - gets message ID
func (m *MessageData) GetID() (string, bool) {
	if m.hasID {
		return m.id, true
	}
	return "", false
}

// ClearID - clears message ID
func (m *MessageData) ClearID() {
	m.hasID = false
	m.id = ""
}

// CopyID - copies message ID
func (m *MessageData) CopyID(m2 *MessageData) {
	m2id, ok := m2.GetID()
	if ok {
		m.hasID = true
		m.id = m2id
	} else {
		m.hasID = false
		m.id = ""
	}
}

// MessageQuit - see protocol.md (quit)
type MessageQuit struct {
	md      *MessageData
	Message string
}

// GetType - MessageQuit.
func (m *MessageQuit) GetType() string {
	return "quit"
}

// GetData - gets MessageData.
func (m *MessageQuit) GetData() *MessageData {
	if m.md == nil {
		m.md = &MessageData{}
	}
	return m.md
}

// MessageFatal - see protocol.md (fatal)
type MessageFatal struct {
	md      *MessageData
	Message string
}

// GetType - MessageFatal.
func (m *MessageFatal) GetType() string {
	return "fatal"
}

// GetData - gets MessageData.
func (m *MessageFatal) GetData() *MessageData {
	if m.md == nil {
		m.md = &MessageData{}
	}
	return m.md
}

// MessageNote - see protocol.md (note)
type MessageNote struct {
	md      *MessageData
	Message string
}

// GetType - MessageNote.
func (m *MessageNote) GetType() string {
	return "note"
}

// GetData - gets MessageData.
func (m *MessageNote) GetData() *MessageData {
	if m.md == nil {
		m.md = &MessageData{}
	}
	return m.md
}

// MessageOk - see protocol.md (ok)
type MessageOk struct {
	md      *MessageData
	Message string
}

// GetType - MessageOk.
func (m *MessageOk) GetType() string {
	return "ok"
}

// GetData - gets MessageData.
func (m *MessageOk) GetData() *MessageData {
	if m.md == nil {
		m.md = &MessageData{}
	}
	return m.md
}

// MessageFail - see protocol.md (fail)
type MessageFail struct {
	md      *MessageData
	Message string
}

// GetType - MessageFail.
func (m *MessageFail) GetType() string {
	return "fail"
}

// GetData - gets MessageData.
func (m *MessageFail) GetData() *MessageData {
	if m.md == nil {
		m.md = &MessageData{}
	}
	return m.md
}

// MessageRegister - see protocol.md (register)
type MessageRegister struct {
	md   *MessageData
	Name string
	Pass string
}

// GetType - MessageRegister.
func (m *MessageRegister) GetType() string {
	return "register"
}

// GetData - gets MessageData.
func (m *MessageRegister) GetData() *MessageData {
	if m.md == nil {
		m.md = &MessageData{}
	}
	return m.md
}

// MessageLogin - see protocol.md (login)
type MessageLogin struct {
	md   *MessageData
	Name string
	Pass string
}

// GetType - MessageLogin.
func (m *MessageLogin) GetType() string {
	return "login"
}

// GetData - gets MessageData.
func (m *MessageLogin) GetData() *MessageData {
	if m.md == nil {
		m.md = &MessageData{}
	}
	return m.md
}

// MessageNewChannel - creates a new channel
type MessageNewChannel struct {
	md    *MessageData
	Name  string
	Topic string
}

// GetType - MessageNewChannel.
func (m *MessageNewChannel) GetType() string {
	return "new-channel"
}

// GetData - gets MessageData.
func (m *MessageNewChannel) GetData() *MessageData {
	if m.md == nil {
		m.md = &MessageData{}
	}
	return m.md
}

// MessageDeleteChannel - deletes a channel
type MessageDeleteChannel struct {
	md   *MessageData
	Name string
}

// GetType - MessageDeleteChannel.
func (m *MessageDeleteChannel) GetType() string {
	return "delete-channel"
}

// GetData - gets MessageData.
func (m *MessageDeleteChannel) GetData() *MessageData {
	if m.md == nil {
		m.md = &MessageData{}
	}
	return m.md
}

// MessageChannelTopic - changes a channel's topic
type MessageChannelTopic struct {
	md    *MessageData
	Name  string
	Topic string
}

// GetType - MessageChannelTopic.
func (m *MessageChannelTopic) GetType() string {
	return "channel-topic"
}

// GetData - gets MessageData.
func (m *MessageChannelTopic) GetData() *MessageData {
	if m.md == nil {
		m.md = &MessageData{}
	}
	return m.md
}

// MessageSubscribe - subscribes to a channel
type MessageSubscribe struct {
	md   *MessageData
	Name string
	//Id string // todo id channel parsing
	Subbed bool
}

// GetType - MessageSubscribe.
func (m *MessageSubscribe) GetType() string {
	return "subscribe"
}

// GetData - gets MessageData.
func (m *MessageSubscribe) GetData() *MessageData {
	if m.md == nil {
		m.md = &MessageData{}
	}
	return m.md
}

// MessagePostMsg - sends a message to a channel
type MessagePostMsg struct {
	md        *MessageData
	Content   string
	Channel   string
	Author    string
	HasAuthor bool
}

// GetType - MessagePostMsg.
func (m *MessagePostMsg) GetType() string {
	return "post-message"
}

// GetData - gets MessageData.
func (m *MessagePostMsg) GetData() *MessageData {
	if m.md == nil {
		m.md = &MessageData{}
	}
	return m.md
}

// MessageGetMsgs - gets messages from the server
type MessageGetMsgs struct {
	md        *MessageData
	ChannelID int
	MessageID int
	Amount    int
}

// GetType - MessageGetMsgs.
func (m *MessageGetMsgs) GetType() string {
	return "get-messages"
}

// GetData - gets MessageData.
func (m *MessageGetMsgs) GetData() *MessageData {
	if m.md == nil {
		m.md = &MessageData{}
	}
	return m.md
}

// MessageGetMsgsResult - sends fetched messages
type MessageGetMsgsResult struct {
	md       *MessageData
	Messages []*ChatMessage
}

// GetType - MessageGetMsgsResult.
func (m *MessageGetMsgsResult) GetType() string {
	return "get-messages-result"
}

// GetData - gets MessageData.
func (m *MessageGetMsgsResult) GetData() *MessageData {
	if m.md == nil {
		m.md = &MessageData{}
	}
	return m.md
}

// MessageListChannels - lists channels
type MessageListChannels struct {
	md          *MessageData
	Channels    map[string]interface{}
	HasChannels bool
}

// GetType - MessageListChannels.
func (m *MessageListChannels) GetType() string {
	return "list-channels"
}

// GetData - gets MessageData.
func (m *MessageListChannels) GetData() *MessageData {
	if m.md == nil {
		m.md = &MessageData{}
	}
	return m.md
}

// MessageListUsers - lists users
type MessageListUsers struct {
	md       *MessageData
	Users    []*UserStatus // todo custom status messages
	HasUsers bool
}

// GetType - MessageListUsers.
func (m *MessageListUsers) GetType() string {
	return "list-users"
}

// GetData - gets MessageData.
func (m *MessageListUsers) GetData() *MessageData {
	if m.md == nil {
		m.md = &MessageData{}
	}
	return m.md
}

// MessageUserQuit - informs clients about someone closing the connection
type MessageUserQuit struct {
	md       *MessageData
	Username string
}

// GetType - MessageUserQuit.
func (m *MessageUserQuit) GetType() string {
	return "user-quit"
}

// GetData - gets MessageData.
func (m *MessageUserQuit) GetData() *MessageData {
	if m.md == nil {
		m.md = &MessageData{}
	}
	return m.md
}

// MessageKick - kicks or kickbans a user
type MessageKick struct {
	md       *MessageData
	ID       int
	Username string
	Ban      bool
}

// GetType - MessageKick.
func (m *MessageKick) GetType() string {
	return "kick"
}

// GetData - gets MessageData.
func (m *MessageKick) GetData() *MessageData {
	if m.md == nil {
		m.md = &MessageData{}
	}
	return m.md
}

// MessageUnregister - indicates a kick/ban event
// todo make the user unregister by his own will
type MessageUnregister struct {
	md       *MessageData
	ID       int
	Username string
}

// GetType - MessageUnregister.
func (m *MessageUnregister) GetType() string {
	return "unregister"
}

// GetData - gets MessageData.
func (m *MessageUnregister) GetData() *MessageData {
	if m.md == nil {
		m.md = &MessageData{}
	}
	return m.md
}

// LoadMessage - builds a MessageBase struct based on given map[string]interface{}
func LoadMessage(iface map[string]interface{}) (BaseMessage, error) {
	var msg BaseMessage

	switch iface["type"].(string) {
	case "quit":
		if _, ok := iface["message"]; !ok {
			return nil, errors.New("no message field in quit message")
		}
		msg = &MessageQuit{Message: iface["message"].(string)}
	case "fatal":
		if _, ok := iface["message"]; !ok {
			return nil, errors.New("no message field in fatal message")
		}
		msg = &MessageFatal{Message: iface["message"].(string)}
	case "note":
		if _, ok := iface["message"]; !ok {
			return nil, errors.New("no message field in note message")
		}
		msg = &MessageNote{Message: iface["message"].(string)}
	case "ok":
		if _, ok := iface["message"]; !ok {
			return nil, errors.New("no message field in ok message")
		}
		msg = &MessageOk{Message: iface["message"].(string)}
	case "fail":
		if _, ok := iface["message"]; !ok {
			return nil, errors.New("no message field in fail message")
		}
		msg = &MessageFail{Message: iface["message"].(string)}
	case "register":
		if _, ok := iface["name"]; !ok {
			return nil, errors.New("no name field in register message")
		}
		if _, ok := iface["pass"]; !ok {
			return nil, errors.New("no pass field in register message")
		}
		msg = &MessageRegister{Name: iface["name"].(string), Pass: iface["pass"].(string)}
	case "login":
		if _, ok := iface["name"]; !ok {
			return nil, errors.New("no name field in login message")
		}
		if _, ok := iface["pass"]; !ok {
			return nil, errors.New("no pass field in login message")
		}
		msg = &MessageLogin{Name: iface["name"].(string), Pass: iface["pass"].(string)}
	case "new-channel":
		if _, ok := iface["name"]; !ok {
			return nil, errors.New("no name field in new-channel message")
		}
		if _, ok := iface["topic"]; !ok {
			return nil, errors.New("no topic field in new-channel message")
		}
		msg = &MessageNewChannel{Name: iface["name"].(string), Topic: iface["topic"].(string)}
	case "delete-channel":
		if _, ok := iface["name"]; !ok {
			return nil, errors.New("no name field in new-channel message")
		}
		msg = &MessageDeleteChannel{Name: iface["name"].(string)}
	case "channel-topic":
		if _, ok := iface["name"]; !ok {
			return nil, errors.New("no name field in channel-topic message")
		}
		if _, ok := iface["topic"]; !ok {
			return nil, errors.New("no topic field in channel-topic message")
		}
		msg = &MessageChannelTopic{Name: iface["name"].(string), Topic: iface["topic"].(string)}
	case "subscribe":
		if _, ok := iface["name"]; !ok {
			return nil, errors.New("no name field in subscribe message")
		}
		if _, ok := iface["subbed"]; !ok {
			return nil, errors.New("no subbed field in subscribe message")
		}
		msg = &MessageSubscribe{Name: iface["name"].(string), Subbed: iface["subbed"].(bool)}
	case "post-message":
		if _, ok := iface["content"]; !ok {
			return nil, errors.New("no content field in send-message message")
		}
		if _, ok := iface["channel"]; !ok {
			return nil, errors.New("no channel field in send-message message")
		}
		author := ""
		hasAuthor := false
		if _, ok := iface["author"]; ok {
			author = iface["author"].(string)
			hasAuthor = true
		}
		msg = &MessagePostMsg{Content: iface["content"].(string), Channel: iface["channel"].(string), Author: author, HasAuthor: hasAuthor}
	case "get-messages":
		if _, ok := iface["channel-id"]; !ok {
			return nil, errors.New("no channel-id field in get-messages message")
		}
		if _, ok := iface["message-id"]; !ok {
			return nil, errors.New("no message-id field in get-messages message")
		}
		if _, ok := iface["amount"]; !ok {
			return nil, errors.New("no amount field in get-messages message")
		}
		msg = &MessageGetMsgs{ChannelID: int(iface["channel-id"].(float64)), MessageID: int(iface["message-id"].(float64)), Amount: int(iface["amount"].(float64))}
	case "get-messages-result":
		if _, ok := iface["messages"]; !ok {
			return nil, errors.New("no messages field in get-messages-result message")
		}
		msg = &MessageGetMsgsResult{Messages: iface["messages"].([]*ChatMessage)}
	case "list-channels":
		channels := map[string]interface{}{}
		hasChannels := false
		if _, ok := iface["channels"]; ok {
			channels = iface["channels"].(map[string]interface{})
			hasChannels = true
		}
		msg = &MessageListChannels{Channels: channels, HasChannels: hasChannels}
	case "list-users":
		users := []*UserStatus{}
		hasUsers := false
		if _, ok := iface["users"]; ok {
			users = iface["users"].([]*UserStatus)
			hasUsers = true
		}
		msg = &MessageListUsers{Users: users, HasUsers: hasUsers}
	case "user-quit":
		if _, ok := iface["username"]; !ok {
			return nil, errors.New("no username field in user-quit message")
		}
	case "kick":
		var hasID bool
		var hasUsername bool
		if _, ok := iface["id"]; ok {
			hasID = true
		}
		if _, ok := iface["username"]; ok {
			hasUsername = true
		}
		if hasID && hasUsername {
			return nil, errors.New("you can't have id and username fields together in kick message")
		} else if !hasID && !hasUsername {
			return nil, errors.New("no id or username field in kick message")
		}
		if _, ok := iface["ban"]; ok {
			if hasID {
				msg = &MessageKick{ID: int(iface["id"].(float64)), Ban: iface["ban"].(bool)}
			} else if hasUsername {
				msg = &MessageKick{Username: iface["username"].(string), Ban: iface["ban"].(bool)}
			}
		} else {
			return nil, errors.New("no ban field in kick message")
		}
	case "unregister":
		if _, ok := iface["id"]; !ok {
			return nil, errors.New("no id field in unregister message")
		}
		if _, ok := iface["username"]; !ok {
			return nil, errors.New("no username field in unregister message")
		}
		msg = &MessageUnregister{ID: int(iface["id"].(float64)), Username: iface["username"].(string)}
	}

	if msg != nil {
		id, ok := iface["_id"].(string)
		if ok {
			l := len(id)
			if l >= 64 {
				l = 63
			}
			msg.GetData().SetID(string(id[0:l]))
		}
		return msg, nil
	}
	return nil, errors.New("invalid type " + iface["type"].(string))
}

// MessageToIface - converts given message to a map[string]interface{}
func MessageToIface(msg BaseMessage) (map[string]interface{}, error) {
	var out map[string]interface{}

	switch msg.(type) {
	case *MessageQuit:
		out = map[string]interface{}{"type": "quit", "message": msg.(*MessageQuit).Message}
	case *MessageFatal:
		out = map[string]interface{}{"type": "fatal", "message": msg.(*MessageFatal).Message}
	case *MessageNote:
		out = map[string]interface{}{"type": "note", "message": msg.(*MessageNote).Message}
	case *MessageOk:
		out = map[string]interface{}{"type": "ok", "message": msg.(*MessageOk).Message}
	case *MessageFail:
		out = map[string]interface{}{"type": "fail", "message": msg.(*MessageFail).Message}
	case *MessageRegister:
		out = map[string]interface{}{"type": "register", "name": msg.(*MessageRegister).Name, "pass": msg.(*MessageRegister).Pass}
	case *MessageLogin:
		out = map[string]interface{}{"type": "login", "name": msg.(*MessageLogin).Name, "pass": msg.(*MessageLogin).Pass}
	case *MessageNewChannel:
		out = map[string]interface{}{"type": "new-channel", "name": msg.(*MessageNewChannel).Name, "topic": msg.(*MessageNewChannel).Topic}
	case *MessageDeleteChannel:
		out = map[string]interface{}{"type": "delete-channel", "name": msg.(*MessageDeleteChannel).Name}
	case *MessageChannelTopic:
		out = map[string]interface{}{"type": "channel-topic", "name": msg.(*MessageChannelTopic).Name, "topic": msg.(*MessageChannelTopic).Topic}
	case *MessageSubscribe:
		out = map[string]interface{}{"type": "subscribe", "name": msg.(*MessageSubscribe).Name, "subbed": msg.(*MessageSubscribe).Subbed}
	case *MessagePostMsg:
		if msg.(*MessagePostMsg).HasAuthor {
			out = map[string]interface{}{"type": "post-message", "content": msg.(*MessagePostMsg).Content, "channel": msg.(*MessagePostMsg).Channel, "author": msg.(*MessagePostMsg).Author}
		} else {
			out = map[string]interface{}{"type": "post-message", "content": msg.(*MessagePostMsg).Content, "channel": msg.(*MessagePostMsg).Channel}
		}
	case *MessageGetMsgs:
		out = map[string]interface{}{"type": "get-messages", "channel-id": msg.(*MessageGetMsgs).ChannelID, "message-id": msg.(*MessageGetMsgs).MessageID, "amount": msg.(*MessageGetMsgs).Amount}
	case *MessageGetMsgsResult:
		out = map[string]interface{}{"type": "get-messages-result", "messages": msg.(*MessageGetMsgsResult).Messages}
	case *MessageListChannels:
		if msg.(*MessageListChannels).HasChannels {
			out = map[string]interface{}{"type": "list-channels", "channels": msg.(*MessageListChannels).Channels}
		} else {
			out = map[string]interface{}{"type": "list-channels"}
		}
	case *MessageListUsers:
		if msg.(*MessageListUsers).HasUsers {
			out = map[string]interface{}{"type": "list-users", "users": msg.(*MessageListUsers).Users}
		} else {
			out = map[string]interface{}{"type": "list-users"}
		}
	case *MessageUserQuit:
		out = map[string]interface{}{"type": "user-quit", "username": msg.(*MessageUserQuit).Username}
	case *MessageKick:
		// todo, it isn't meant to be sent by server anyways
	default:
		return nil, errors.New("invalid type")
	}

	id, ok := msg.GetData().GetID()
	if ok {
		out["_id"] = id
	}

	return out, nil
}
