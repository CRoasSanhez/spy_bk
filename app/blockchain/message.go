package blockchain

import (
	"bytes"
)

// Message struct
// Options
// Data is the bytes to be sent
type Message struct {
	Identifier byte
	Options    []byte
	Data       []byte
	Checksum   uint32

	Reply chan Message
}

func Initialize(id byte) *Message {
	return &Message{Identifier: id}
}

func (m *Message) MarshallBinary() ([]byte, error) {
	buff := new(bytes.Buffer)
	buff.Write(FitBytesInto(m.Options, MESSAGE_OPTIONS_SIZE))
	buff.Write(m.Data)
	return buff.Bytes(), nil
}

// UnmarshallBinary
func (m *Message) UnmarshallBinary(data []byte) error {

	return nil
}
