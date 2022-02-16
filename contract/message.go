package contract

type Message struct {
	Data []byte
}

func NewMessage(data []byte) Message {
	return Message{
		Data: data,
	}
}

func (m Message) String() string {
	return string(m.Data)
}
