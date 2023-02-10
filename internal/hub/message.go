package hub

import (
	"encoding/json"
	"github.com/hashicorp/memberlist"
)

type Message struct {
	Payload string
	notify  chan<- struct{}
}

func (m *Message) Invalidates(other memberlist.Broadcast) bool {
	return false
}
func (m *Message) Finished() {
	if m.notify != nil {
		close(m.notify)
	}
}
func (m *Message) Message() []byte {
	data, err := json.Marshal(m)
	if err != nil {
		return []byte("")
	}
	return data
}

func ParseMessage(data []byte) (*Message, bool) {
	msg := &Message{}
	if err := json.Unmarshal(data, &msg); err != nil {
		return nil, false
	}
	return msg, true
}
