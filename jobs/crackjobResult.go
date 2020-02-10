package jobs

import (
	"bytes"
	"encoding/gob"

	uuid "github.com/satori/go.uuid"
)

// CrackJobResult bundles a client's result of attempting to crack something.
type CrackJobResult struct {
	// Success determines whether the client thinks it succeeded
	Success bool
	// The Payload contains a string representation of the password cracked by the client
	Payload string
	// JobID is a reference to the original job
	JobID uuid.UUID
}

func (result *CrackJobResult) Encode() ([]byte, error) {
	buffer := new(bytes.Buffer)
	encoder := gob.NewEncoder(buffer)
	err := encoder.Encode(result)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func DecodeResult(data []byte) (*CrackJobResult, error) {
	tmp := bytes.NewBuffer(data)
	tmpStruct := new(CrackJobResult)
	decoder := gob.NewDecoder(tmp)
	err := decoder.Decode(tmpStruct)
	if err != nil {
		return nil, err
	}

	return tmpStruct, nil
}
