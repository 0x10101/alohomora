package jobs

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"math/big"
	"time"

	uuid "github.com/satori/go.uuid"
)

// A CrackJob is something the server sends to a client in order for that
// client to crack it.
type CrackJob struct {
	// The Type of cracking the client should attempt
	Type JobType
	// ID is a unique ID to be able to identify the job later on
	ID uuid.UUID
	// The Payload contains the actual thing to crack (handshake, hash, ...)
	Payload []byte
	// Gen contains the parameters for generating passwords
	Gen *GenerationParams
	// Started is a timestamp that is set when the job is transmitted to a client (used for timeouts)
	Started time.Time
}

func (job *CrackJob) String() string {
	return fmt.Sprintf("%s", job.ID.String()[:8])
}

func DecodeJob(data []byte) (*CrackJob, error) {
	tmp := bytes.NewBuffer(data)
	tmpStruct := new(CrackJob)
	decoder := gob.NewDecoder(tmp)
	err := decoder.Decode(tmpStruct)
	if err != nil {
		return nil, err
	}

	return tmpStruct, nil

}

func (job *CrackJob) DecodeWPA2() (*WPA2Payload, error) {
	tmp := bytes.NewBuffer(job.Payload)
	tmpStruct := new(WPA2Payload)
	decoder := gob.NewDecoder(tmp)
	err := decoder.Decode(tmpStruct)
	if err != nil {
		return nil, err
	}
	return tmpStruct, nil
}

func (job *CrackJob) Encode() ([]byte, error) {
	buffer := new(bytes.Buffer)
	encoder := gob.NewEncoder(buffer)
	err := encoder.Encode(job)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func NewWPA2Job(data []byte, charset []rune, length int64, offset *big.Int, amount int64, essid, bssid string) (*CrackJob, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}

	params := &GenerationParams{Charset: charset, Length: length, Offset: offset, Amount: amount}
	wpa2Payload := &WPA2Payload{Data: data, ESSID: essid, BSSID: bssid}
	payload, err := wpa2Payload.Encode()

	if err != nil {
		return nil, err
	}

	return &CrackJob{
		Type:    WPA2,
		ID:      id,
		Payload: payload,
		Gen:     params,
	}, nil
}
