package globaloption

import (
	"errors"
	"strings"
)

const OptionNameIgnoreLabels string = "ignore-labels"

type IgnoreLabel string

func parseIgnoreLabel(v string) (IgnoreLabel, error) {
	trimmed := strings.TrimSpace(v)
	if trimmed == "" {
		return "", errors.New("empty ignore label is not allowed")
	}
	return IgnoreLabel(trimmed), nil
}

func (i IgnoreLabel) String() string {
	return string(i)
}

type IgnoreLabels []IgnoreLabel

func ParseIgnoreLabels(v string) (IgnoreLabels, error) {
	if v == "" {
		return IgnoreLabels{}, nil
	}

	values := strings.Split(v, ",")
	ignoreLabels := make([]IgnoreLabel, 0, len(values))
	for _, value := range values {
		ignoreLabel, err := parseIgnoreLabel(value)
		if err != nil {
			return nil, err
		}

		ignoreLabels = append(ignoreLabels, ignoreLabel)
	}

	return ignoreLabels, nil
}

func (i IgnoreLabels) Strings() []string {
	result := make([]string, 0, len(i))
	for _, label := range i {
		result = append(result, label.String())
	}
	return result
}
