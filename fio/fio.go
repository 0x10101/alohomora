package fio

import "os"

// WriteTo is a convenience function that writes a string (data) to a
// file identified by its path.
// Any error occurring while doing so will be returned.
// The file will be closed upon return.
func WriteTo(filepath, data string) error {
	f, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(data)
	if err != nil {
		return err
	}

	return nil
}
