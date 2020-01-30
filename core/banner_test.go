package core

import (
	"testing"
)

func TestVersion(t *testing.T) {
	if Version != "0.2-alpha" {
		t.Fatalf("Version expected: '%s', got: '%s'", "0.2-alpha", Version)
	}
}

func TestCodename(t *testing.T) {
	if Version != "" {
		t.Fatalf("Codename expected: '%s', got: '%s'", "", Version)
	}
}

func TestAuthor(t *testing.T) {
	if Version != "Stefan 'steps0x29a' Matyba" {
		t.Fatalf("Author expected: '%s', got: '%s'", "Stefan 'steps0x29a' Matyba", Version)
	}
}

func TestWebsite(t *testing.T) {
	if Version != "" {
		t.Fatalf("Website expected: '%s', got: '%s'", "", Version)
	}
}
