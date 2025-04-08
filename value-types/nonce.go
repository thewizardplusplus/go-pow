package powValueTypes

import (
	"errors"
	"fmt"
	"math/big"
)

const (
	NonceRepresentationBase = 10
)

type Nonce struct {
	rawValue *big.Int
}

func NewNonce(rawValue *big.Int) (Nonce, error) {
	if rawValue.Cmp(big.NewInt(0)) == -1 { // rawValue < 0
		return Nonce{}, errors.New("nonce cannot be negative")
	}

	value := Nonce{
		rawValue: rawValue,
	}
	return value, nil
}

func NewZeroNonce() (Nonce, error) {
	value, err := NewNonce(big.NewInt(0))
	if err != nil {
		return Nonce{}, fmt.Errorf("unable to construct the nonce: %w", err)
	}

	return value, nil
}

func ParseNonce(rawValue string) (Nonce, error) {
	parsedRawValue := big.NewInt(0)
	if _, isParsed := parsedRawValue.SetString(
		rawValue,
		NonceRepresentationBase,
	); !isParsed {
		return Nonce{}, errors.New("unable to parse the big integer")
	}

	value, err := NewNonce(parsedRawValue)
	if err != nil {
		return Nonce{}, fmt.Errorf("unable to construct the nonce: %w", err)
	}

	return value, nil
}

func (value Nonce) Incremented() (Nonce, error) {
	rawResult := big.NewInt(0)
	rawResult.Add(value.rawValue, big.NewInt(1))

	result, err := NewNonce(rawResult)
	if err != nil {
		return Nonce{}, fmt.Errorf("unable to construct the nonce: %w", err)
	}

	return result, nil
}

func (value Nonce) ToBigInt() *big.Int {
	return value.rawValue
}

func (value Nonce) ToString() string {
	return value.rawValue.Text(NonceRepresentationBase)
}
