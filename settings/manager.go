package settings

import (
	"fmt"
	"strings"
	"sync"

	providers "github.com/foxesknow/tango/settings/internal"
)

var allSettings = make(map[string]Provider)
var lock sync.Mutex

func init() {
	allSettings["env"] = providers.NewEnvironmentSettings()
	allSettings["os"] = providers.NewOSSettings()
}

// Checks to see if a provider is registered
func IsRegistered(provider string) bool {
	_, present := getProvider(normalizeName(provider))
	return present
}

// Returns the value for a setting.
// The setting must be in Provider:Name format
func Value(setting string) (string, bool) {
	provider, name, err := extractProviderAndName(setting)
	if err != nil {
		return "", false
	}

	p, ok := getProvider(normalizeName(provider))
	if !ok {
		return "", false
	}

	return p.GetSetting(name)
}

func normalizeName(name string) string {
	return strings.ToLower(name)
}

func getProvider(name string) (provider Provider, found bool) {
	lock.Lock()
	defer lock.Unlock()

	provider, found = allSettings[name]
	return
}

func extractProviderAndName(value string) (string, string, error) {
	pivot := strings.Index(value, ":")

	if pivot == -1 {
		return "", "", fmt.Errorf("could not find :")
	}

	provider := value[:pivot]
	setting := value[pivot+1:]

	return provider, setting, nil
}