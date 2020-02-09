package gen

import (
	"errors"
	"math/big"
	"testing"

	"github.com/steps0x29a/alohomora/bigint"
)

func TestGenerate(t *testing.T) {

	numbers := []rune("0123456789")

	var table = []struct {
		c   []rune
		l   int64
		s   *big.Int
		pw  string
		err error
	}{
		{numbers, 3, bigint.ToBigInt("123"), "123", nil},
		{numbers, 3, bigint.ToBigInt("999"), "999", nil},
		{numbers, 3, bigint.ToBigInt("1000"), "", errors.New("Amount is too large")},
		{numbers, 64, bigint.ToBigInt("12312312390"), "0000000000000000000000000000000000000000000000000000012312312390", nil},
	}

	for _, tt := range table {
		pw, err := GeneratePassword(tt.c, tt.l, tt.s)
		if pw != tt.pw {
			t.Errorf("Expected '%s', got '%s' from (%s, %d, %s)", tt.pw, pw, string(tt.c), tt.l, tt.s)
		}

		if err != nil && tt.err == nil {
			t.Errorf("Expected nil error but got '%s' from (%s, %d, %s)", err, string(tt.c), tt.l, tt.s)
		}

		if err == nil && tt.err != nil {
			t.Errorf("Expected error '%s' but got nil error from (%s, %d, %s)", tt.err, string(tt.c), tt.l, tt.s)
		}
	}

}
