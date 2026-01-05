package runtimeoption

import "fmt"

type Mode string

const ModeTwoApprove Mode = "two-approve"

func ParseMode(v string) (Mode, error) {
	if v == "" {
		return "", fmt.Errorf("--mode option is required")
	}
	mode := Mode(v)

	//exhaustive:enforce
	//nolint:gocritic
	switch mode {
	case ModeTwoApprove:
		return mode, nil
	}

	return "", fmt.Errorf("invalid mode. mode: %s", mode)
}

func (m Mode) String() string {
	return string(m)
}
