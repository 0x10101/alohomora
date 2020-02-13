package jobs

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"math/big"
	"time"

	"github.com/steps0x29a/alohomora/gen"
	"github.com/steps0x29a/alohomora/term"

	uuid "github.com/satori/go.uuid"
	"github.com/steps0x29a/alohomora/bigint"
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

// The CrackJobInfo type is used to report information on a CrackJob via the server's
// REST API as a JSON object.
type CrackJobInfo struct {
	Type    string    `json:"type"`
	ID      string    `json:"id"`
	Started time.Time `json:"started"`
	Charset string    `json:"charset"`
	Length  int64     `json:"length"`
	Amount  int64     `json:"amount"`
	Offset  *big.Int  `json:"offset"`
	First   string    `json:"first"`
	Last    string    `json:"last"`
}

// String calculates and returns a CrackJob's short ID and returns it as a string.
// A CrackJob's short ID are the first 8 characters of its UUID.
func (job *CrackJob) String() string {
	return fmt.Sprintf("%s", job.ID.String()[:8])
}

// Info builds and returns a CrackJobInfo object from a CrackJob object in order
// to send it via the server's REST API (as JSON).
func (job *CrackJob) Info() *CrackJobInfo {

	first, err := gen.GeneratePassword(job.Gen.Charset, job.Gen.Length, job.Gen.Offset)
	if err != nil {
		term.Warn("Unable to calculate first password of job %s: %s\n", job.ShortID(), err)
		first = "<err>"
	}

	amnt := big.NewInt(job.Gen.Amount)
	last, err := gen.GeneratePassword(job.Gen.Charset, job.Gen.Length, bigint.Sub(bigint.Add(job.Gen.Offset, amnt), big.NewInt(1)))
	if err != nil {
		term.Warn("Unable to calculate last password of job %s: %s\n", job.ShortID(), err)
		last = "<err>"
	}

	return &CrackJobInfo{
		Type:    job.Type.String(),
		ID:      job.ID.String()[:8],
		Started: job.Started,
		Charset: string(job.Gen.Charset),
		Length:  job.Gen.Length,
		Amount:  job.Gen.Amount,
		Offset:  bigint.Cp(job.Gen.Offset),
		First:   first,
		Last:    last,
	}
}

// DecodeJob decodes a CrackJob from a slice of bytes.
// If decoding the raw bytes succeeds, the newly decoded object is returned.
// In case an error occurrs, that error is returned instead.
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

// ShortID returns the first eight characters of a job's ID as a string.
func (job *CrackJob) ShortID() string {
	return job.ID.String()[:8]
}

// DecodeWPA2 attempts to decode a CrackJob's payload as a WPA2Payload.
// If decoding succeeds, the decoded WPA2Payload is returned. If it fails, an error
// is returned instead.
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

// Encode encodes a CrackJob object to a byte slice that can be sent via
// a socket connection.
func (job *CrackJob) Encode() ([]byte, error) {
	buffer := new(bytes.Buffer)
	encoder := gob.NewEncoder(buffer)
	err := encoder.Encode(job)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

// NewWPA2Job generates a new WPA2 crack job from the given parameters.
// The function requires the handshake's raw bytes (data), a password length (length),
// and offset for password generation (offset), the amount of passwords to generate for
// attempting this job (amount), an ESSID and a BSSID.
// If everything works as expected, the newly generated CrackJob is returned.
// In case of an error, that error is returned instead.
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
