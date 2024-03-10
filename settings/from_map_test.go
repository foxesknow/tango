package settings

import (
	"testing"
)

func Test_FromMap_Create_HasSetting(t *testing.T) {
	provider := createProvider()

	value, found := provider.GetSetting("Name")
	if !found {
		t.Fatal("could not find Name")
	}

	if value != "Jack" {
		t.Fatal("Expected Jack")
	}
}

func Test_FromMap_Create_NoSetting(t *testing.T) {
	provider := createProvider()

	value, found := provider.GetSetting("Location")
	if found {
		t.Fatal("Location should not exist!")
	}

	if value != "" {
		t.Error("Expected an empty string")
	}
}

func createProvider() Provider {
	values := make(map[string]string)
	values["Name"] = "Jack"
	values["Age"] = "48"

	return FromMap(values)
}
