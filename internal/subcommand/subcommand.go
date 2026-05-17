package subcommand

import (
	"fmt"
)

type Name string

const (
	Help       Name = "help"
	Version    Name = "version"
	TwoApprove Name = "two-approve"
)

func ParseName(v string) (Name, error) {
	switch Name(v) {
	case Help:
		return Help, nil
	case Version:
		return Version, nil
	case TwoApprove:
		return TwoApprove, nil
	}

	return "", fmt.Errorf("invalid subcommand. subcommand: %s", v)
}

func (n Name) String() string {
	return string(n)
}
