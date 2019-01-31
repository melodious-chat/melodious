package main

func handler(mel *Melodious, messageStream chan<- BaseMessage, message BaseMessage) {
	switch message.(type) {
	case *MessageRegister:
		messageStream <- &MessageFail{Message: "not implemented"}
	case *MessageLogin:
		messageStream <- &MessageFail{Message: "not implemented"}
		// todo
	}
}

// handler wrapper for REST (might be removed later)
func wrapHandler(mel *Melodious, messageStream <-chan BaseMessage, f func(mel *Melodious, messageStream <-chan BaseMessage, message BaseMessage)) func(BaseMessage) {
	return func(message BaseMessage) {
		f(mel, messageStream, message)
	}
}
