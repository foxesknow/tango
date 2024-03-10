package settings

import (
	"testing"
)

func Test_Cons(t *testing.T) {
	head := NameValue("Name", "Jack")
	tail := NameValue("Location", "Island")
	provider := Cons(head, tail)

	value, found := provider.GetSetting("Name")
	if !found {
		t.Fatal("Expected to find Name setting")
	}

	if value != "Jack" {
		t.Fatal("Expected value to be Jack")
	}

	value, found = provider.GetSetting("Location")
	if !found {
		t.Fatal("Expected to find Location setting")
	}

	if value != "Island" {
		t.Fatal("Expected value to be Island")
	}
}
