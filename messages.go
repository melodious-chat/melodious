package main

// BaseMessage - A base struct for all messages
type BaseMessage struct {
	Type string `json:"type"`
}

// MessageFatal - see protocol.md (fatal)
type MessageFatal struct {
	BaseMessage
	Message string `json:"message"`
}

// MessageNote - see protocol.md (note)
type MessageNote struct {
	BaseMessage
	Message string `json:"message"`
}

// MessageOk - see protocol.md (ok)
type MessageOk struct {
	BaseMessage
	Message string `json:"message"`
}

// MessageFail - see protocol.md (fail)
type MessageFail struct {
	BaseMessage
	Message string `json:"message"`
}

// MessageRegister - see protocol.md (register)
type MessageRegister struct {
	BaseMessage
	Name string `json:"name"`
	Pass string `json:"pass"`
}

// MessageLogin - see protocol.md (login)
type MessageLogin struct {
	BaseMessage
	Name string `json:"name"`
	Pass string `json:"pass"`
}

// MessagePing - see protocol.md (ping)
type MessagePing struct {
	BaseMessage
}

// MessagePong - see protocol.md (pong)
type MessagePong struct {
	BaseMessage
}
