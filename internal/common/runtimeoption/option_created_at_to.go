package runtimeoption

import (
	"fmt"
	"time"
)

type CreatedAtTo time.Time

func ParseCreatedAtTo(v string) (CreatedAtTo, error) {
	if v == "" {
		return CreatedAtTo{}, fmt.Errorf("--to option is required")
	}

	date, err := time.ParseInLocation("20060102", v, time.Local)
	if err != nil {
		return CreatedAtTo{}, fmt.Errorf("invalid input date. value: %s, err: %w", v, err)
	}

	return CreatedAtTo(date), nil
}

func (c CreatedAtTo) Time() time.Time {
	return time.Time(c)
}
