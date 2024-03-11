package text

import (
	"errors"
	"fmt"
	"strings"

	"github.com/foxesknow/tango/settings"
)

type ExpandTextConfig struct {
	Begin                 string
	End                   string
	UnknownVariableLookup func(string) (any, bool)
}

type formatter interface {
	Format(string) string
}

func ExpandText(value string, config ExpandTextConfig) (string, error) {
	if len(config.Begin) == 0 {
		return "", errors.New("ExpandText: begin not set")
	}

	if len(config.End) == 0 {
		return "", errors.New("ExpandText: end not set")
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
	variable, formatting, action, err := extractTokenParts(token)
	if err != nil {
		return "", err
	}

	namespace := ""
	originalVariable := variable

	pivot := strings.Index(variable, settings.NamespaceVariableSeparator)
	if pivot != -1 {
		namespace = variable[0:pivot]
		variable = variable[pivot+1:]
	}

	var lookedUpValue any
	var found bool

	if namespace == "" {
		if config.UnknownVariableLookup != nil {
			lookedUpValue, found = config.UnknownVariableLookup(originalVariable)
			if !found {
				err = fmt.Errorf("UnknownTokenLookup could not resolve '%s'", originalVariable)
			}
		} else {
			err = fmt.Errorf("no UnknownVariableLookup set to resolve '%s'", originalVariable)
		}
	} else {
		provider := settings.GetProvider(namespace)
		if provider != nil {
			lookedUpValue, found = provider.GetSetting(variable)
			if !found {
				err = fmt.Errorf("could not get setting: '%s'", variable)
			}
		} else {
			if config.UnknownVariableLookup != nil {
				lookedUpValue, found = config.UnknownVariableLookup(originalVariable)
				if !found {
					err = fmt.Errorf("UnknownTokenLookup could not resolve '%s'", originalVariable)
				}
			} else {
				err = fmt.Errorf("could not resolve namespace '%s'", namespace)
			}
		}
	}

	if err != nil {
		value = ""
		return
	}

	value, err = applyFormatting(lookedUpValue, formatting)
	if err == nil {
		value, err = applyAction(value, action)
	}

	return
}

func applyFormatting(value any, format string) (string, error) {
	if format == "" {
		return fmt.Sprint(value), nil
	}

	if i := value.(formatter); i != nil {
		return i.Format(format), nil
	}

	return "", errors.New("value does not support formatting")
}

func applyAction(value string, action string) (string, error) {
	switch strings.ToLower(action) {
	case "":
		return value, nil
	case "tolower":
		return strings.ToLower(value), nil
	case "toupper":
		return strings.ToUpper(value), nil
	case "trim":
		return strings.TrimSpace(value), nil
	default:
		return "", fmt.Errorf("unsupported action: '%s'", action)
	}
}

func extractTokenParts(token string) (variable string, formatting string, action string, err error) {
	parts := strings.Split(token, "|")

	switch len(parts) {
	case 1:
		variable = parts[0]
		formatting = ""
		action = ""
		err = nil
	case 2:
		variable = parts[0]
		formatting = parts[1]
		action = ""
		err = nil
	case 3:
		variable = parts[0]
		formatting = parts[1]
		action = parts[2]
		err = nil
	default:
		variable = ""
		formatting = ""
		err = fmt.Errorf("too many parts in token: %d", len(parts))
	}

	return
}
