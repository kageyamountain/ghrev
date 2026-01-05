package runtimeoption

import "fmt"

type RepositoryOwner string

func ParseRepositoryOwner(v string) (RepositoryOwner, error) {
	if v == "" {
		return "", fmt.Errorf("--owner option is required")
	}

	return RepositoryOwner(v), nil
}

func (r RepositoryOwner) String() string {
	return string(r)
}
