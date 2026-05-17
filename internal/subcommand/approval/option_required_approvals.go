package approval

import "fmt"

const OptionNameRequiredApprovals string = "required-approvals"

type RequiredApprovals int

func ParseRequiredApprovals(v int) (RequiredApprovals, error) {
	if v <= 0 {
		return 0, fmt.Errorf("--required-approvals option is required and must be >= 1")
	}

	return RequiredApprovals(v), nil
}

func (r RequiredApprovals) Int() int {
	return int(r)
}
