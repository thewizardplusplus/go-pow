package powValueTypes

import (
	"errors"
	"fmt"
	"time"
)

const (
	CreatedAtRepresentationFormat = time.RFC3339Nano
)

type CreatedAt struct {
	rawValue time.Time
}

func NewCreatedAt(rawValue time.Time) (CreatedAt, error) {
	if rawValue.IsZero() {
		return CreatedAt{}, errors.New("`CreatedAt` timestamp cannot be zero time")
	}

	value := CreatedAt{
		rawValue: rawValue,
	}
	return value, nil
}

func ParseCreatedAt(rawValue string) (CreatedAt, error) {
	parsedRawValue, err := time.Parse(CreatedAtRepresentationFormat, rawValue)
	if err != nil {
		return CreatedAt{}, fmt.Errorf("unable to parse the time: %w", err)
	}

	value, err := NewCreatedAt(parsedRawValue)
	if err != nil {
		return CreatedAt{}, fmt.Errorf(
			"unable to construct the `CreatedAt` timestamp: %w",
			err,
		)
	}

	return value, nil
}

func (value CreatedAt) ToTime() time.Time {
	return value.rawValue
}

func (value CreatedAt) ToString() string {
	return value.rawValue.Format(CreatedAtRepresentationFormat)
}
