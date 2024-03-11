package providers

import "os"

type EnvironmentSettings struct {
}

func NewEnvironmentSettings() *EnvironmentSettings {
	return &EnvironmentSettings{}
}

func (self *EnvironmentSettings) GetSetting(name string) (value any, found bool) {
	value, found = os.LookupEnv(name)
	return
}
