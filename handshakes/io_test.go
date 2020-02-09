package handshakes

import (
	"errors"
	"testing"
)

func TestNewHandshake(t *testing.T) {
	hs := NewHandshake()
	if hs.BSSID != "" {
		t.Errorf("Expected empty string in empty handshake BSSID, but got '%s'", hs.BSSID)
	}

	if hs.ESSID != "" {
		t.Errorf("Expected empty string in empty handshake ESSID, but got '%s'", hs.ESSID)
	}

	if hs.Data != nil {
		t.Errorf("Expected nil byte array as empty handshake data, but got '%s'", hs.Data)
	}

}

func TestParseFilename(t *testing.T) {

	defErr := errors.New("Filename does not match requirements")

	var table = []struct {
		in    string
		bssid string
		essid string
		err   error
	}{
		{"SOMEBSSID_001122334455.pcap", "SOMEBSSID", "00:11:22:33:44:55", nil},
		{"Completely_wrong.docx", "", "", defErr},
	}

	for _, tt := range table {
		bssid, essid, err := parseHandshakeFilename(tt.in)
		if essid != tt.essid {
			t.Errorf("got '%s', expected '%s' from '%s'", essid, tt.essid, tt.in)
		}

		if bssid != tt.bssid {
			t.Errorf("got '%s', expected '%s' from '%s'", bssid, tt.bssid, tt.in)
		}

		if err != nil && tt.err == nil {
			t.Errorf("got '%s' but expected nil error from '%s'", err, tt.in)
		}

		if err == nil && tt.err != nil {
			t.Errorf("got nil error, but expected '%s' from '%s'", tt.err, tt.in)
		}
	}
}
