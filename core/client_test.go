package core

import (
	"testing"

	uuid "github.com/satori/go.uuid"
)

func TestNewClient(t *testing.T) {

	c := newClient(nil)

	if c.ID != uuid.Nil {
		t.Fatalf("Expected nil ID for new client, got: '%s'", c.ID.String())
	}

}

func TestShortID(t *testing.T) {
	c := newClient(nil)

	sID := c.ShortID()
	if sID != "00000000" {
		t.Fatalf("Expected client short ID: '%s', got: '%s'", "00000000", c.ShortID())
	}
}
