package runtimeoption

import "fmt"

type RepositoryName string

func ParseRepositoryName(v string) (RepositoryName, error) {
	if v == "" {
		return "", fmt.Errorf("--name option is required")
	}

	return RepositoryName(v), nil
}

func (r RepositoryName) String() string {
	return string(r)
}
