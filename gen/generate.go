package gen

import (
	"errors"
	"math/big"

	"github.com/steps0x29a/alohomora/bigint"
)

const (
	emptyRune rune = '\U0010FFFF'
)

// GeneratePassword generates a password from a given charset with a given length.
// In order to do so, it will shift the characters in the charset by the amount
// provided by shift.
// For example, if charset is 0123456789, length is 3 and shift is 13, it will generate
// 013. In other words: given the three values, the generation of passwords for
// bruteforcing can be distributed to any number of workers.
func GeneratePassword(charset []rune, length int64, shift *big.Int) (string, error) {
	var buffer = make([]rune, length)
	var size = big.NewInt(int64(len(charset)))
	bigLength := big.NewInt(length)

	// The max value is the maximum rotation possible with the given input data
	var maxValue = bigint.Sub(bigint.Pow(size, bigLength), big.NewInt(1))

	if bigint.Lt(maxValue, shift) {
		// maxValue is smaller than requested amount, error!
		return "", errors.New("Amount is too large")
	}

	// Empty out the buffer
	var i int64
	for i = 0; i < length; i++ {
		buffer[i] = emptyRune
	}

	// Allocate 'pointers' (automatically zeroed-out)
	var pointers = make([]int, length)

	// remaining rotation value, modified in loop below
	var remaining = bigint.Cp(shift)

	// Maximum rotation value for any given position (at first index, that is)
	var maxRotVal = bigint.Pow(size, bigint.Sub(bigLength, big.NewInt(1)))

	zero := big.NewInt(0)

	// Loop until remaining rotations are 0
	for bigint.Gt(remaining, zero) {
		// while remaining is larger than rotation value of first position...
		for bigint.Gt(remaining, maxRotVal) {
			// Rotate first rune
			pointers[0]++
			remaining = bigint.Sub(remaining, maxRotVal)
		}
		// iterate over all positions
		for i := length - 1; i >= 0; i-- {
			// Get rotational value for current position
			var rotVal = bigint.Pow(size, bigint.Sub(bigint.Sub(bigLength, big.NewInt(1)), big.NewInt(i)))
			// If rotational value is larger than remaning rotations
			if bigint.Gt(rotVal, remaining) {
				// find position to rotate (one to the right)
				var rotPos = bigint.Add(big.NewInt(i), big.NewInt(1))
				// Make sure we don't get out of range
				if bigint.Gt(rotPos, bigint.Sub(bigLength, big.NewInt(1))) {
					rotPos.Set(bigint.Sub(bigLength, big.NewInt(1)))
				}
				// Calculate rotational value
				rotVal = bigint.Pow(size, bigint.Sub(bigint.Sub(bigLength, big.NewInt(1)), rotPos))
				remaining = bigint.Sub(remaining, rotVal)
				// Rotate position by one
				pointers[rotPos.Int64()]++
				break

			} else if bigint.Eq(rotVal, remaining) {
				// Remaining is same as current position's rotational value, simply rotate
				pointers[i]++
				remaining = bigint.Sub(remaining, rotVal)
				break
			}
		}
	}

	// Populate buffer with runes from charset according to their current rotation
	for i, p := range pointers {
		buffer[i] = charset[p]
	}

	// Make and return string
	return string(buffer), nil
}
