package bytes

// EndsWith tests whether a given []byte ends with another given []byte.
// Returns true if the last bytes in buffer equal the bytes in value
func EndsWith(buffer, value []byte) bool {
	if len(buffer) < len(value) {
		return false
	}
	fromIndex := len(buffer) - len(value)

	sub := buffer[fromIndex:]

	for i, b := range sub {
		if b != value[i] {
			return false
		}
	}

	return true
}
