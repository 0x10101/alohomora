package jobs

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

type HashPayload struct {
	Data []byte
	Salt []byte
}

// Encode encodes a HashPayload object to an array of bytes.
// The encoded bytes are returned as well as any errors.
func (payload *HashPayload) Encode() ([]byte, error) {
	buffer := new(bytes.Buffer)
	encoder := gob.NewEncoder(buffer)
	err := encoder.Encode(payload)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func (payload *HashPayload) String() string {
	if len(payload.Salt) > 0 {
		return fmt.Sprintf("%s:%s", string(payload.Data), string(payload.Salt))
	} else {
		return string(payload.Data)
	}

}
