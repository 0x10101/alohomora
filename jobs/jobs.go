package jobs

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"math/big"
	"time"

	uuid "github.com/satori/go.uuid"
)

// JobType is an enum that defines the type of crackjob (WPA2, MD5, ...)
type JobType int

const (
	// WPA2 identifies a WPA2 crackjob
	WPA2 JobType = 1
)

// A WPA2Payload contains the handshake data as well as an ESSID and a BSSID.
// This will change in the future as clients will be parsing PCAP files on
// their own, extracting the required information themselves.
type WPA2Payload struct {
	// Data contains the raw capture data (PCAP)
	Data []byte
	// ESSID of target
	ESSID string
	// BSSID of target
	BSSID string
}

// CrackJobResult bundles a client's result of attempting to crack something.
type CrackJobResult struct {
	// Success determines whether the client thinks it succeeded
	Success bool
	// The Payload contains a string representation of the password cracked by the client
	Payload string
	// JobID is a reference to the original job
	JobID uuid.UUID
}

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

func (params *GenerationParams) Encode() ([]byte, error) {
	buffer := new(bytes.Buffer)
	encoder := gob.NewEncoder(buffer)
	err := encoder.Encode(params)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
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

func (payload *WPA2Payload) Encode() ([]byte, error) {
	buffer := new(bytes.Buffer)
	encoder := gob.NewEncoder(buffer)
	err := encoder.Encode(payload)
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
