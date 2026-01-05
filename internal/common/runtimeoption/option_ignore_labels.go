package runtimeoption

import (
	"strings"
)

type IgnoreLabel string

type IgnoreLabels []IgnoreLabel

func ParseIgnoreLabel(v string) ([]string, error) {
	if v == "" {
		return []string{}, nil
	}

	values := strings.Split(v, ",")
	ignoreLabels := make([]string, len(values))
	for _, value := range values {
		trimmed := strings.TrimSpace(value)
		ignoreLabels = append(ignoreLabels, trimmed)
	}

	return ignoreLabels, nil
}
