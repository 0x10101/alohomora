package core

import (
	"testing"
)

func TestVersion(t *testing.T) {
	if Version != "0.3" {
		t.Fatalf("Version expected: '%s', got: '%s'", "0.2-alpha", Version)
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
	if Website != "" {
		t.Fatalf("Website expected: '%s', got: '%s'", "", Website)
	}
}
