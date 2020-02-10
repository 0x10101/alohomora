package fio

import (
	"os"
	"path/filepath"

	uuid "github.com/satori/go.uuid"
)

// TempFilePath builds and returns a temporary file path.
// The temporary file will have a UUID as its name and be located in
// os.TempDir().
// If generating the UUID fails, an error is returned.
func TempFilePath() (string, error) {
	uid, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	return filepath.Join(os.TempDir(), uid.String()), nil
}
