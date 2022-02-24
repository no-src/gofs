package contract

// Message the file sync message
type Message struct {
	Data []byte
}

// NewMessage create an instance of Message
func NewMessage(data []byte) Message {
	return Message{
		Data: data,
	}
}

// String convert the message data to string
func (m Message) String() string {
	return string(m.Data)
}
