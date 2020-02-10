package ext

import (
	"testing"
)

func TestEnvKey(t *testing.T) {
	if aircrackExecutableKey != "AIRCRACK" {
		t.Fatalf("Expected: '%s', got '%s'", "AIRCRACK", aircrackExecutableKey)
	}
}

func TestKeyFromOutput(t *testing.T) {
	test := make(map[string]string)

	test["[Something]"] = "Something"
	test["Nothing"] = ""
	test["[some[evilthing]]"] = "some[evilthing]"
	test[""] = ""
	test["[  blah   ]"] = "blah"

	for key, value := range test {
		val := KeyFromOutput(key)
		if val != value {
			t.Fatalf("Expected: '%s' from '%s', got '%s'", value, key, val)
		}
	}
}
