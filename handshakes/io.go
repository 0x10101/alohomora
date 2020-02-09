package handshakes

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/steps0x29a/alohomora/term"
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

// Read reads a file identified by path and tries to parse it as a WPA2 handshake.
// The read data will then be saved to the Handshake object the function is called
// on. If the handshake contains the ESSID and BSSID of the AP, those will be
// extracted. If they can't be found, the function attempts to extract that info
// from the PCAP's filename (see parseHandshakeFilename for details). If that also
// fails, an error will be returned.
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

// parseHandshakeFilenam takes a filename (not a path!) and attempts to extract a BSSID and ESSID from it.
// In order for this to work, the file should be named <ESSID>_<BSSID>.pcap and the BSSID should not contain
// any divider characters (ABCDEF0123456 instead of AB:CD:EF:12:34:56). For printing purposes, the function
// will splice in colons after every two characters (ABCDEF0123456 becomes AB:CD:EF:12:34:56) on its own.
// This function returns the extracted BSSID and ESSID (in that order) or empty strings and an error.
func parseHandshakeFilename(filename string) (essid string, bssid string, err error) {

	// We use a compiled regex here
	re := regexp.MustCompile(`^(.*)_([0-9a-zA-Z]{12}(\..*))$`)

	// If the filename does not match, don't even attempt parsing it
	match := re.MatchString(filename)
	if !match {
		return "", "", errors.New("Filename does not match requirements")
	}

	// Split on the underscore character, remove .pcap extension
	split := strings.Split(strings.Replace(filename, ".pcap", "", 1), "_")

	// Splice in colons and return
	return split[0], term.Splice(split[1], ':', 2), nil
}
