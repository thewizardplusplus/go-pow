package powValueTypes

import (
	"errors"
)

type TargetBitIndex struct {
	rawValue int
}

func NewTargetBitIndex(rawValue int) (TargetBitIndex, error) {
	if rawValue < 0 {
		return TargetBitIndex{}, errors.New("target bit index cannot be negative")
	}

	value := TargetBitIndex{
		rawValue: rawValue,
	}
	return value, nil
}

func (value TargetBitIndex) ToInt() int {
	return value.rawValue
}
