package opts

import (
	"github.com/steps0x29a/islazy/term"
)

// CleanupCharset removes duplicate runes from a given charset.
func CleanupCharset(charset []rune) ([]rune, error) {
	term.Warn("Charset cleanup is not yet implemented\n")
	return charset, nil

}
