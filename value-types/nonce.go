package powValueTypes

import (
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"math/big"
)

const (
	NonceRepresentationBase = 10
)

type Nonce struct {
	rawValue *big.Int
}

func NewNonce(rawValue *big.Int) (Nonce, error) {
	if rawValue.Sign() < 0 {
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

type RandomNonceParams struct {
	RandomReader io.Reader
	MinRawValue  *big.Int
	MaxRawValue  *big.Int
}

func NewRandomNonce(params RandomNonceParams) (Nonce, error) {
	rawValueRange := big.NewInt(0).Sub(params.MaxRawValue, params.MinRawValue)
	if rawValueRange.Sign() < 0 {
		return Nonce{}, errors.New("raw value range cannot be negative")
	}
	if rawValueRange.Sign() == 0 {
		return Nonce{}, errors.New("raw value range cannot be zero")
	}

	randomRawValue, err := rand.Int(params.RandomReader, rawValueRange)
	if err != nil {
		return Nonce{}, fmt.Errorf(
			"unable to generate the random big integer: %w",
			err,
		)
	}

	randomValue, err := NewNonce(
		big.NewInt(0).Add(randomRawValue, params.MinRawValue),
	)
	if err != nil {
		return Nonce{}, fmt.Errorf("unable to construct the nonce: %w", err)
	}

	return randomValue, nil
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
