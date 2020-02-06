package core

import (
	"testing"
)

func TestVersion(t *testing.T) {
	if Version != "0.4" {
		t.Fatalf("Version expected: '%s', got: '%s'", "0.4", Version)
	}
}

func TestCodename(t *testing.T) {
	if Codename != "" {
		t.Fatalf("Codename expected: '%s', got: '%s'", "", Codename)
	}
}

func TestAuthor(t *testing.T) {
	if Author != "Stefan 'steps0x29a' Matyba" {
		t.Fatalf("Author expected: '%s', got: '%s'", "Stefan 'steps0x29a' Matyba", Codename)
	}
}

func TestWebsite(t *testing.T) {
	if Website != "https://github.com/steps0x29a/alohomora" {
		t.Fatalf("Website expected: '%s', got: '%s'", "https://github.com/steps0x29a/alohomora", Website)
	}
}
