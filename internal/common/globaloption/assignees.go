package globaloption

import (
	"errors"
	"strings"
)

const OptionNameAssignees string = "assignees"

type Assignee string

func parseAssignee(v string) (Assignee, error) {
	trimmed := strings.TrimSpace(v)
	if trimmed == "" {
		return "", errors.New("empty assignee is not allowed")
	}
	return Assignee(trimmed), nil
}

func (a Assignee) String() string {
	return string(a)
}

type Assignees []Assignee

func ParseAssignees(v string) (Assignees, error) {
	v = strings.Trim(v, ", ")
	if v == "" {
		return Assignees{}, nil
	}

	values := strings.Split(v, ",")
	assignees := make([]Assignee, 0, len(values))
	for _, value := range values {
		assignee, err := parseAssignee(value)
		if err != nil {
			return nil, err
		}

		assignees = append(assignees, assignee)
	}

	return assignees, nil
}

func (a Assignees) Strings() []string {
	result := make([]string, 0, len(a))
	for _, assignee := range a {
		result = append(result, assignee.String())
	}
	return result
}
