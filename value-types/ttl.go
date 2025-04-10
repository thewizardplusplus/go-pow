package powValueTypes

import (
	"errors"
	"fmt"
	"time"
)

type TTL struct {
	rawValue time.Duration
}

func NewTTL(rawValue time.Duration) (TTL, error) {
	if rawValue < 0 {
		return TTL{}, errors.New("TTL cannot be negative")
	}

	value := TTL{
		rawValue: rawValue,
	}
	return value, nil
}

func ParseTTL(rawValue string) (TTL, error) {
	parsedRawValue, err := time.ParseDuration(rawValue)
	if err != nil {
		return TTL{}, fmt.Errorf("unable to parse the duration: %w", err)
	}

	value, err := NewTTL(parsedRawValue)
	if err != nil {
		return TTL{}, fmt.Errorf("unable to construct the TTL: %w", err)
	}

	return value, nil
}

func (value TTL) ToDuration() time.Duration {
	return value.rawValue
}

func (value TTL) ToString() string {
	return value.rawValue.String()
}
