package providers

import (
	"testing"
)

func TestUser(t *testing.T) {
	s := EnvironmentSettings{}

	value, found := s.GetSetting("USER")

	if !found {
		t.Fatal("user should have been found")
	}

	if value == "" {
		t.Fatal("Got an empty value instead of a name")
	}
}

func TestUserWrongCase(t *testing.T) {
	s := EnvironmentSettings{}

	_, found := s.GetSetting("user")

	if found {
		t.Fatal("user should have been found")
	}
}
