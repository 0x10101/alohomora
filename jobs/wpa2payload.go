package jobs

import (
	"bytes"
	"encoding/gob"
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

func (payload *WPA2Payload) Encode() ([]byte, error) {
	buffer := new(bytes.Buffer)
	encoder := gob.NewEncoder(buffer)
	err := encoder.Encode(payload)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}
