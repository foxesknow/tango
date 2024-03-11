package providers

import (
	"fmt"
	"os"
	"strings"
)

type OSSettings struct {
	values map[string]string
}

func NewOSSettings() *OSSettings {
	settings := OSSettings{
		values: make(map[string]string),
	}

	if value, err := os.UserHomeDir(); err == nil {
		settings.values["homedir"] = value
	}

	if value, err := os.Hostname(); err == nil {
		settings.values["hostname"] = value
	}

	if value, err := os.UserHomeDir(); err == nil {
		settings.values["homedir"] = value
	}

	settings.values["pid"] = fmt.Sprintf("%d", os.Getpid())
	settings.values["tempdir"] = os.TempDir()

	return &settings
}

func (self *OSSettings) GetSetting(name string) (value any, found bool) {
	value, found = self.values[strings.ToLower(name)]
	return
}
