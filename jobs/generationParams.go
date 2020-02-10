package jobs

import (
	"bytes"
	"encoding/gob"
	"math/big"
)

// GenerationParams bundle everything a client needs to know in order to
// generate passwords for bruteforcing something.
type GenerationParams struct {
	// Charset is the character set the passwords should be generated from
	Charset []rune
	// Length is the length of the passwords that the client should generate
	Length int64
	// Offset is the starting point for password generation
	Offset *big.Int
	// Amount is - obviously - the amount of passwords the client should generate
	Amount int64
}

func (params *GenerationParams) Encode() ([]byte, error) {
	buffer := new(bytes.Buffer)
	encoder := gob.NewEncoder(buffer)
	err := encoder.Encode(params)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}
