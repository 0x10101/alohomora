package jobs

import (
	"testing"
)

func TestString(t *testing.T) {
	var table = []struct {
		t JobType
		s string
	}{
		{WPA2, "WPA2"},
		{MD5, "MD5"},
	}

	for _, x := range table {
		ts := x.t.String()
		if ts != x.s {
			t.Errorf("Expected '%s', got '%s' from (%d)", x.s, ts, x.t)
		}
	}

}
