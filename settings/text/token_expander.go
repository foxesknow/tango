package text

import (
	"errors"
	"strings"

	"github.com/foxesknow/tango/settings"
)

type ExpandTextConfig struct {
	Begin              string
	End                string
	UnknownTokenLookup func(string) (string, error)
}

func ExpandText(value string, config ExpandTextConfig) (string, error) {
	if len(config.Begin) == 0 {
		return "", errors.New("Begin not set")
	}

	if len(config.End) == 0 {
		return "", errors.New("End not set")
	}

	beginLen := len(config.Begin)
	endLen := len(config.End)

	var builder strings.Builder

	for {
		index := strings.Index(value, config.Begin)

		if index == -1 {
			break
		}

		builder.WriteString(value[0:index])
		value = value[index+beginLen:]

		endIndex := strings.Index(value, config.End)
		if endIndex == -1 {
			break
		}

		token := value[0:endIndex]
		expandedToken, err := expandToken(token, config)
		if err != nil {
			return "", err
		}

		builder.WriteString(expandedToken)

		value = value[endIndex+endLen:]
	}

	builder.WriteString(value)

	return builder.String(), nil
}

func expandToken(token string, config ExpandTextConfig) (value string, err error) {
	variable, formatting, err := extractTokenParts(token)
	if err != nil {
		return "", err
	}

	namespace := ""
	originalVariable := variable

	pivot := strings.Index(variable, ":")
	if pivot != -1 {
		namespace = variable[0:pivot]
		variable = variable[pivot+1:]
	}

	if namespace == "" {
		if config.UnknownTokenLookup != nil {
			value, err = config.UnknownTokenLookup(originalVariable)
		}
	} else {
		provider := settings.GetProvider(namespace)
		if provider != nil {
			settingValue, found := provider.GetSetting(variable)
			if found {
				value = settingValue
			} else {
				err = errors.New("could not get setting")
			}
		} else {
			if config.UnknownTokenLookup != nil {
				value, err = config.UnknownTokenLookup(originalVariable)
			}
		}
	}

	if err != nil {
		value = ""
		return
	}

	value, err = applyFormatting(value, formatting)
	return
}

func applyFormatting(value string, formatting string) (string, error) {
	switch strings.ToLower(formatting) {
	case "tolower":
		return strings.ToLower(value), nil
	case "toupper":
		return strings.ToUpper(value), nil
	case "trim":
		return strings.TrimSpace(value), nil
	default:
		return "", errors.New("unsupported formatting")
	}
}

func extractTokenParts(token string) (variable string, formatting string, err error) {
	parts := strings.Split(token, "|")

	switch len(parts) {
	case 1:
		variable = parts[0]
		formatting = ""
		err = nil
	case 2:
		variable = parts[0]
		formatting = parts[2]
		err = nil
	default:
		variable = ""
		formatting = ""
		err = errors.New("too many parts in token")
	}

	return
}
