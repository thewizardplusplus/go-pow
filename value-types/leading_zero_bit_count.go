package powValueTypes

import (
	"errors"
)

type LeadingZeroBitCount struct {
	rawValue int
}

func NewLeadingZeroBitCount(rawValue int) (LeadingZeroBitCount, error) {
	if rawValue < 0 {
		return LeadingZeroBitCount{}, errors.New(
			"leading zero bit count cannot be negative",
		)
	}

	value := LeadingZeroBitCount{
		rawValue: rawValue,
	}
	return value, nil
}

func (value LeadingZeroBitCount) ToInt() int {
	return value.rawValue
}
