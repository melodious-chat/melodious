package main

// messageHandler - handles messages received from users
func messageHandler(mel *Melodious, messageStream chan<- BaseMessage, message BaseMessage) {
	switch message.(type) {
	case *MessageRegister:
		messageStream <- &MessageFail{Message: "not implemented"}
	case *MessageLogin:
		messageStream <- &MessageFail{Message: "not implemented"}
		// todo
	}
}

// wrapMessageHandler - wraps a message handler to allow passing it without explicitly passing some context-specific data
func wrapMessageHandler(
	mel *Melodious, messageStream <-chan BaseMessage, f func(mel *Melodious, messageStream <-chan BaseMessage, message BaseMessage)) func(BaseMessage) {
	return func(message BaseMessage) {
		f(mel, messageStream, message)
	}
}
