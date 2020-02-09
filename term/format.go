package term

import "bytes"

// Splice inserts a given rune into a given string at a given interval and returns the modified string.
func Splice(raw string, insert rune, interval int) string {
	if interval <= 0 || interval >= len(raw) {
		return raw
	}
	var buffer bytes.Buffer

	last := len(raw) - 1
	for i, char := range raw {
		buffer.WriteRune(char)
		if i%interval == interval-1 && i != last {
			buffer.WriteRune(insert)
		}
	}
	return buffer.String()
}
