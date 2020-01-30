package handshakes

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/steps0x29a/islazy/term"
)

// Handshake bundles everything we need to know about a handshake
type Handshake struct {
	ESSID string
	BSSID string
	Data  []byte
}

// NewHandshake creates a new instance of the Handshake type.
func NewHandshake() *Handshake {
	return new(Handshake)
}

func (handshake *Handshake) Read(path string) error {
	filename := filepath.Base(path)

	essid, bssid, err := HandshakeInfo(path)
	if err != nil {
		// Attempt parsing the handshake filename as a last resort
		essid, bssid, err = parseHandshakeFilename(filename)
		if err != nil {
			return err
		}
	}

	if essid == "" {
		return errors.New("Unable to determine ESSID, can't crack this handshake")
	}

	if bssid == "" {
		return errors.New("Unable to determine BSSID, can't crack this handshake")
	}

	f, err := os.Open(path)
	if err != nil {
		return err
	}

	defer f.Close()

	data := make([]byte, 0)
	buffer := make([]byte, 1024*1024*1024)

	for {
		read, err := f.Read(buffer)
		if read == 0 {
			break
		}
		if err != nil && err != io.EOF {
			return err
		}

		sub := buffer[:read]
		data = append(data, sub...)
	}

	handshake.ESSID = essid
	handshake.BSSID = bssid
	handshake.Data = data
	return nil
}

func parseHandshakeFilename(filename string) (essid string, bssid string, err error) {
	// Filename is supposed to be <ESSID>_<BSSID>.pcap
	split := strings.Split(strings.Replace(filename, ".pcap", "", 1), "_")
	if len(split) != 2 {
		return "", "", errors.New("Filename does not match requirements")
	}

	return split[0], term.InsertAfterEvery(split[1], ':', 2), nil
}
