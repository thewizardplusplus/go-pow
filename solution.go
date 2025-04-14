package pow

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/samber/mo"
	powValueTypes "github.com/thewizardplusplus/go-pow/value-types"
)

type Solution struct {
	challenge Challenge
	nonce     powValueTypes.Nonce
	hashSum   mo.Option[powValueTypes.HashSum]
}

func (entity Solution) Challenge() Challenge {
	return entity.challenge
}

func (entity Solution) Nonce() powValueTypes.Nonce {
	return entity.nonce
}

func (entity Solution) HashSum() mo.Option[powValueTypes.HashSum] {
	return entity.hashSum
}

func (entity Solution) Verify() error {
	hashData, err := entity.challenge.hashDataLayout.Execute(ChallengeHashData{
		Challenge: entity.challenge,
		Nonce:     entity.nonce,
	})
	if err != nil {
		return fmt.Errorf("unable to execute the hash data layout: %w", err)
	}

	hashSum := entity.challenge.hash.ApplyTo(hashData)

	expectedHashSum, isPresent := entity.hashSum.Get()
	if isPresent && !bytes.Equal(hashSum.ToBytes(), expectedHashSum.ToBytes()) {
		return errors.New("hash sum doesn't match the expected one")
	}

	targetBitIndex, err := entity.challenge.TargetBitIndex()
	if err != nil {
		return fmt.Errorf("unable to get the target bit index: %w", err)
	}

	target := makeTarget(targetBitIndex)
	if !isHashSumFitTarget(hashSum, target) {
		return errors.New("hash sum doesn't fit the target")
	}

	return nil
}
