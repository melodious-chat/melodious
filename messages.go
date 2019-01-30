package main

import "errors"

// BaseMessage - A base struct for all messages
type BaseMessage interface {
	GetType() string
	Validate() bool
}

// MessageQuit - see protocol.md (quit)
type MessageQuit struct {
	Message string
}

// GetType - MessageQuit.
func (m *MessageQuit) GetType() string {
	return "quit"
}

// Validate - MessageQuit.
func (m *MessageQuit) Validate() bool {
	return true
}

// MessageFatal - see protocol.md (fatal)
type MessageFatal struct {
	Message string
}

// GetType - MessageFatal.
func (m *MessageFatal) GetType() string {
	return "fatal"
}

// Validate - MessageFatal.
func (m *MessageFatal) Validate() bool {
	return true
}

// MessageNote - see protocol.md (note)
type MessageNote struct {
	Message string
}

// GetType - MessageNote.
func (m *MessageNote) GetType() string {
	return "note"
}

// Validate - MessageNote.
func (m *MessageNote) Validate() bool {
	return true
}

// MessageOk - see protocol.md (ok)
type MessageOk struct {
	Message string
}

// GetType - MessageOk.
func (m *MessageOk) GetType() string {
	return "ok"
}

// Validate - MessageOk.
func (m *MessageOk) Validate() bool {
	return true
}

// MessageFail - see protocol.md (fail)
type MessageFail struct {
	Message string
}

// GetType - MessageFail.
func (m *MessageFail) GetType() string {
	return "fail"
}

// Validate - MessageFail.
func (m *MessageFail) Validate() bool {
	return true
}

// MessageRegister - see protocol.md (register)
type MessageRegister struct {
	Name string
	Pass string
}

// GetType - MessageRegister.
func (m *MessageRegister) GetType() string {
	return "register"
}

// Validate - MessageRegister.
func (m *MessageRegister) Validate() bool {
	return true
}

// MessageLogin - see protocol.md (login)
type MessageLogin struct {
	Name string
	Pass string
}

// GetType - MessageLogin.
func (m *MessageLogin) GetType() string {
	return "login"
}

// Validate - MessageLogin.
func (m *MessageLogin) Validate() bool {
	return true
}

// MessagePing - see protocol.md (ping)
type MessagePing struct{}

// GetType - MessagePing.
func (m *MessagePing) GetType() string {
	return "ping"
}

// Validate - MessagePing.
func (m *MessagePing) Validate() bool {
	return true
}

// MessagePong - see protocol.md (pong)
type MessagePong struct{}

// GetType - MessagePong.
func (m *MessagePong) GetType() string {
	return "pong"
}

// Validate - MessagePong.
func (m *MessagePong) Validate() bool {
	return true
}

// LoadMessage - builds a MessageBase struct based on given map[string]interface{}
func LoadMessage(iface map[string]interface{}) (BaseMessage, error) {
	switch iface["type"].(string) {
	case "quit":
		return &MessageQuit{Message: iface["message"].(string)}, nil
	case "fatal":
		return &MessageFatal{Message: iface["message"].(string)}, nil
	case "note":
		return &MessageNote{Message: iface["message"].(string)}, nil
	case "ok":
		return &MessageOk{Message: iface["message"].(string)}, nil
	case "fail":
		return &MessageFail{Message: iface["message"].(string)}, nil
	case "register":
		return &MessageRegister{Name: iface["name"].(string), Pass: iface["pass"].(string)}, nil
	case "login":
		return &MessageLogin{Name: iface["name"].(string), Pass: iface["pass"].(string)}, nil
	case "ping":
		return &MessagePing{}, nil
	case "pong":
		return &MessagePong{}, nil
	}
	return nil, errors.New("invalid type " + iface["type"].(string))
}

// MessageToIface - converts given message to a map[string]interface{}
func MessageToIface(msg BaseMessage) (map[string]interface{}, error) {
	switch msg.(type) {
	case *MessageQuit:
		return map[string]interface{}{"type": "quit", "message": msg.(*MessageQuit).Message}, nil
	case *MessageFatal:
		return map[string]interface{}{"type": "fatal", "message": msg.(*MessageFatal).Message}, nil
	case *MessageNote:
		return map[string]interface{}{"type": "note", "message": msg.(*MessageNote).Message}, nil
	case *MessageOk:
		return map[string]interface{}{"type": "ok", "message": msg.(*MessageOk).Message}, nil
	case *MessageFail:
		return map[string]interface{}{"type": "fail", "message": msg.(*MessageFail).Message}, nil
	case *MessageRegister:
		return map[string]interface{}{"type": "register", "name": msg.(*MessageRegister).Name, "pass": msg.(*MessageRegister).Pass}, nil
	case *MessageLogin:
		return map[string]interface{}{"type": "login", "name": msg.(*MessageLogin).Name, "pass": msg.(*MessageLogin).Pass}, nil
	case *MessagePing:
		return map[string]interface{}{"type": "ping"}, nil
	case *MessagePong:
		return map[string]interface{}{"type": "pong"}, nil
	}
	return nil, errors.New("invalid type")
}
