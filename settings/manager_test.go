package settings

import (
	"testing"
)

func TestIsRegistered(t *testing.T) {
	registered := IsRegistered("env")
	if !registered {
		t.Error("env not registered")
	}
}

func TestScopedName(t *testing.T) {
	value, found := Value("os:homedir")
	if !found {
		t.Fatal("could not resolve setting")
	}

	if value == "" {
		t.Fatal("value is empty")
	}
}
