package pow

import (
	"bytes"
	"errors"
	"fmt"

	powValueTypes "github.com/thewizardplusplus/go-pow/value-types"
)

type Solution struct {
	challenge Challenge
	nonce     powValueTypes.Nonce
	hashSum   powValueTypes.HashSum
}

func (entity Solution) Challenge() Challenge {
	return entity.challenge
}

func (entity Solution) Nonce() powValueTypes.Nonce {
	return entity.nonce
}

func (entity Solution) HashSum() powValueTypes.HashSum {
	return entity.hashSum
}

func (entity Solution) Verify() error {
	targetBitIndex := entity.challenge.TargetBitIndex()

	hashData, err := entity.challenge.hashDataLayout.Execute(ChallengeHashData{
		Challenge:      entity.challenge,
		TargetBitIndex: targetBitIndex,
		Nonce:          entity.nonce,
	})
	if err != nil {
		return fmt.Errorf("unable to execute the hash data layout: %w", err)
	}

	hashSum := entity.challenge.hash.ApplyTo(hashData)
	if !bytes.Equal(hashSum.ToBytes(), entity.hashSum.ToBytes()) {
		return errors.New("hash sums don't match")
	}

	target := makeTarget(targetBitIndex)
	if !isHashSumFitTarget(hashSum, target) {
		return errors.New("hash sum doesn't fit the target")
	}

	return nil
}
