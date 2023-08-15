package message

import (
	"fmt"
)

// New --
func New() *Message {
	return &Message{}
}

// Message represents a unmarshaled message or a message which to be marshaled
type Message struct {
	Type       byte   // message type
	ID         uint   // unique id, zero while notify mode
	Route      string // route for locating service
	Data       []byte // payload
	compressed bool   // is message compressed
}

// String --
func (m *Message) String() string {
	return fmt.Sprintf("Type: %s, ID: %d, Route: %s, Compressed: %t, BodyLength: %d",
		types[m.Type],
		m.ID,
		m.Route,
		m.compressed,
		len(m.Data))
}

// Encode --
func (m *Message) Encode() ([]byte, error) {
	return Encode(m)
}
