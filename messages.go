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
	default:
		return nil, errors.New("invalid type")
	}

	id, ok := msg.GetData().GetID()
	if ok {
		out["_id"] = id
	}

	return out, nil
}
