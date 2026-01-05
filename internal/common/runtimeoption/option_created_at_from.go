package runtimeoption

import (
	"fmt"
	"time"
)

type CreatedAtFrom time.Time

func ParseCreatedAtFrom(v string) (CreatedAtFrom, error) {
	if v == "" {
		return CreatedAtFrom{}, fmt.Errorf("--fron option is required")
	}

	date, err := time.ParseInLocation("20060102", v, time.Local)
	if err != nil {
		return CreatedAtFrom{}, fmt.Errorf("invalid input date. value: %s, err: %w", v, err)
	}

	return CreatedAtFrom(date), nil
}

func (c CreatedAtFrom) Time() time.Time {
	return time.Time(c)
}
