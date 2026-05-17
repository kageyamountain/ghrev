package subcommand

import (
	"fmt"
)

type Name string

const (
	Help     Name = "help"
	Version  Name = "version"
	Approval Name = "approval"
)

func ParseName(v string) (Name, error) {
	switch Name(v) {
	case Help:
		return Help, nil
	case Version:
		return Version, nil
	case Approval:
		return Approval, nil
	}

	return "", fmt.Errorf("invalid subcommand. subcommand: %s", v)
}

func (n Name) String() string {
	return string(n)
}
