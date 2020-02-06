package fio

import "os"

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
