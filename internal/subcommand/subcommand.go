package subcommand

import (
	"fmt"
)

type Name string

const (
	Help        Name = "help"
	Version     Name = "version"
	Approval    Name = "approval"
	FirstReview Name = "first-review"
)

func ParseName(v string) (Name, error) {
	switch Name(v) {
	case Help:
		return Help, nil
	case Version:
		return Version, nil
	case Approval:
		return Approval, nil
	case FirstReview:
		return FirstReview, nil
	}

	return "", fmt.Errorf("invalid subcommand. subcommand: %s", v)
}

func (n Name) String() string {
	return string(n)
}
