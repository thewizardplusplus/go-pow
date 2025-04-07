package powValueTypes

import (
	"errors"
)

type LeadingZeroCount struct {
	rawValue int
}

func NewLeadingZeroCount(rawValue int) (LeadingZeroCount, error) {
	if rawValue < 0 {
		return LeadingZeroCount{}, errors.New("leading zero count cannot be negative")
	}

	value := LeadingZeroCount{
		rawValue: rawValue,
	}
	return value, nil
}

func (value LeadingZeroCount) ToInt() int {
	return value.rawValue
}
