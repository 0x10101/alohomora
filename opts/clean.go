package opts

// CleanupCharset removes duplicate runes from a given charset.
func CleanupCharset(charset []rune) []rune {
	m := make(map[rune]bool)
	n := make([]rune, 0)
	for _, r := range charset {
		m[r] = true
	}

	for r := range m {
		n = append(n, r)
	}

	return n

}
