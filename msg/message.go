package msg

import (
	"bytes"
	"encoding/gob"

	uuid "github.com/satori/go.uuid"
)

// MessageType determines the type of message we are dealing with
type MessageType int

const (
	// None indicates an unindentified message
	None = 0
	// Hello identifies a client hello message
	Hello = 1
	// Ack identifies messages that the server sends to clients as a response to their Hello
	Ack = 2
	// Idle identifies messages claiming that a client is idle and needs work
	Idle = 3
	// Task identifies a message that contains work for a client
	Task = 4
	// Finished identifies a message containing a result (client to server)
	Finished = 5
	// ClientError identifies a message containing a client error
	ClientError = 6
	// Leave identifies a message that asks a client to disconnect
	Leave = 7
)

// A Message is a wrapper around raw bytes sent through a socket connection
// between client and server.
type Message struct {
	Type    MessageType
	ID      uuid.UUID
	Payload []byte
}

// NewMessage creates and initialzes a new message with the given type t and a payload.
func NewMessage(t MessageType, payload []byte) *Message {
	id, err := uuid.NewV4()
	if err != nil {
		id = uuid.Nil
	}

	if payload == nil {
		payload = make([]byte, 0)
	}

	return &Message{
		Type:    t,
		ID:      id,
		Payload: payload,
	}
}

// Decode attempts to decode a given []byte to an instance of Message and
// returns it (or an error)
func Decode(data []byte) (*Message, error) {
	tmp := bytes.NewBuffer(data)
	holder := new(Message)
	decoder := gob.NewDecoder(tmp)
	err := decoder.Decode(holder)
	if err != nil {
		return nil, err
	}
	return holder, nil
}

// Encode takes a message and serializes it to a []byte.
func (message *Message) Encode() ([]byte, error) {
	buffer := new(bytes.Buffer)
	encoder := gob.NewEncoder(buffer)
	err := encoder.Encode(message)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}
