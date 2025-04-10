package pow

import (
	"fmt"
	"time"

	"github.com/samber/mo"
	powValueTypes "github.com/thewizardplusplus/go-pow/value-types"
)

type ChallengeHashData struct {
	Challenge      Challenge
	TargetBitIndex powValueTypes.TargetBitIndex
	Nonce          powValueTypes.Nonce
}

type Challenge struct {
	leadingZeroBitCount powValueTypes.LeadingZeroBitCount
	createdAt           mo.Option[powValueTypes.CreatedAt]
	ttl                 mo.Option[powValueTypes.TTL]
	resource            mo.Option[powValueTypes.Resource]
	serializedPayload   powValueTypes.SerializedPayload
	hash                powValueTypes.Hash
	hashDataLayout      powValueTypes.HashDataLayout
}

func (entity Challenge) LeadingZeroBitCount() powValueTypes.LeadingZeroBitCount { //nolint:lll
	return entity.leadingZeroBitCount
}

func (entity Challenge) TargetBitIndex() (powValueTypes.TargetBitIndex, error) {
	rawValue := entity.hash.SizeInBits() - entity.leadingZeroBitCount.ToInt()

	value, err := powValueTypes.NewTargetBitIndex(rawValue)
	if err != nil {
		return powValueTypes.TargetBitIndex{}, fmt.Errorf(
			"unable to construct the target bit index: %w",
			err,
		)
	}

	return value, nil
}

func (entity Challenge) CreatedAt() mo.Option[powValueTypes.CreatedAt] {
	return entity.createdAt
}

func (entity Challenge) TTL() mo.Option[powValueTypes.TTL] {
	return entity.ttl
}

func (entity Challenge) IsAlive() bool {
	createdAt, isCreatedAtPresent := entity.createdAt.Get()
	ttl, isTTLPresent := entity.ttl.Get()
	return !isCreatedAtPresent ||
		!isTTLPresent ||
		time.Since(createdAt.ToTime()) <= ttl.ToDuration()
}

func (entity Challenge) Resource() mo.Option[powValueTypes.Resource] {
	return entity.resource
}

func (entity Challenge) SerializedPayload() powValueTypes.SerializedPayload {
	return entity.serializedPayload
}

func (entity Challenge) Hash() powValueTypes.Hash {
	return entity.hash
}

func (entity Challenge) HashDataLayout() powValueTypes.HashDataLayout {
	return entity.hashDataLayout
}

func (entity Challenge) Solve() (Solution, error) {
	targetBitIndex, err := entity.TargetBitIndex()
	if err != nil {
		return Solution{}, fmt.Errorf("unable to get the target bit index: %w", err)
	}

	target := makeTarget(targetBitIndex)

	nonce, err := powValueTypes.NewZeroNonce()
	if err != nil {
		return Solution{}, fmt.Errorf(
			"unable to construct the zero initial nonce: %w",
			err,
		)
	}

	var hashSum powValueTypes.HashSum
	for {
		hashData, err := entity.hashDataLayout.Execute(ChallengeHashData{
			Challenge:      entity,
			TargetBitIndex: targetBitIndex,
			Nonce:          nonce,
		})
		if err != nil {
			return Solution{}, fmt.Errorf(
				"unable to execute the hash data layout: %w",
				err,
			)
		}

		hashSum = entity.hash.ApplyTo(hashData)
		if isHashSumFitTarget(hashSum, target) {
			break
		}

		nonce, err = nonce.Incremented()
		if err != nil {
			return Solution{}, fmt.Errorf("unable to increment the nonce: %w", err)
		}
	}

	solution, err := NewSolutionBuilder().
		SetChallenge(entity).
		SetNonce(nonce).
		SetHashSum(hashSum).
		Build()
	if err != nil {
		return Solution{}, fmt.Errorf("unable to build the solution: %w", err)
	}

	return solution, nil
}
