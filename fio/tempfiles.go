package fio

import (
	"os"
	"path/filepath"

	uuid "github.com/satori/go.uuid"
)

func TempFilePath() (string, error) {
	uid, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	return filepath.Join(os.TempDir(), uid.String()), nil
}
