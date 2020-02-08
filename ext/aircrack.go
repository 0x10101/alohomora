package ext

import (
	"os"
	"os/exec"
	"strings"

	"github.com/steps0x29a/alohomora/term"
)

const (
	aircrackExecutableKey string = "AIRCRACK"
)

// AircrackAvailable determines whether or not aircrack-ng is present
// on the machine executing alohomora.
func AircrackAvailable() bool {

	_, err := getAircrackExecutable()
	if err != nil {
		term.Error("Unable to find aircrack-ng: %s\n", err)
		return false
	}

	return true
}

func getAircrackExecutable() (string, error) {
	value, present := os.LookupEnv(aircrackExecutableKey)
	if present {
		return value, nil
	}

	exe, err := exec.LookPath("aircrack-ng")
	if err != nil {
		return "", err
	}
	return exe, nil
}

// Aircrack executes the aircrack-ng command, returning its output as a string
func Aircrack(bssid, essid, wordlist, target string) (string, error) {
	exe, _ := getAircrackExecutable()
	var cmdline = []string{"-q", "-a2", "-b", bssid, "-e", essid, "-w", wordlist, target}
	cmd := exec.Command(exe, cmdline...)

	data, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// KeyFromOutput extracts the key that aircrack-ng has found from its output
func KeyFromOutput(output string) string {
	begin := strings.Index(output, "[") + 1
	end := strings.LastIndex(output, "]")
	if end <= begin {
		return ""
	}

	return strings.TrimSpace(output[begin:end])
}
