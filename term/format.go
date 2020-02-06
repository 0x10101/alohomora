package term

import "bytes"

func InsertAfterEvery(raw string, insert rune, interval int) string {
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

func Reverse(raw string) string {
	c := []rune(raw)
	n := len(c)

	for i := 0; i < n/2; i++ {
		c[i], c[n-1-i] = c[n-1-i], c[i]
	}

	return string(c)
}
