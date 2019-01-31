package main

import (
	"errors"
)

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

// LoadMessage - builds a MessageBase struct based on given map[string]interface{}
func LoadMessage(iface map[string]interface{}) (BaseMessage, error) {
	switch iface["type"].(string) {
	case "quit":
		if _, ok := iface["message"]; !ok {
			return nil, errors.New("no message field in quit message")
		}
		return &MessageQuit{Message: iface["message"].(string)}, nil
	case "fatal":
		if _, ok := iface["message"]; !ok {
			return nil, errors.New("no message field in fatal message")
		}
		return &MessageFatal{Message: iface["message"].(string)}, nil
	case "note":
		if _, ok := iface["message"]; !ok {
			return nil, errors.New("no message field in note message")
		}
		return &MessageNote{Message: iface["message"].(string)}, nil
	case "ok":
		if _, ok := iface["message"]; !ok {
			return nil, errors.New("no message field in ok message")
		}
		return &MessageOk{Message: iface["message"].(string)}, nil
	case "fail":
		if _, ok := iface["message"]; !ok {
			return nil, errors.New("no message field in fail message")
		}
		return &MessageFail{Message: iface["message"].(string)}, nil
	case "register":
		if _, ok := iface["name"]; !ok {
			return nil, errors.New("no name field in register message")
		}
		if _, ok := iface["pass"]; !ok {
			return nil, errors.New("no pass field in register message")
		}
		return &MessageRegister{Name: iface["name"].(string), Pass: iface["pass"].(string)}, nil
	case "login":
		if _, ok := iface["name"]; !ok {
			return nil, errors.New("no name field in login message")
		}
		if _, ok := iface["pass"]; !ok {
			return nil, errors.New("no pass field in login message")
		}
		return &MessageLogin{Name: iface["name"].(string), Pass: iface["pass"].(string)}, nil
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
	}
	return nil, errors.New("invalid type")
}
