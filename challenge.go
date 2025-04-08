package pow

import (
	"fmt"

	"github.com/samber/mo"
	powValueTypes "github.com/thewizardplusplus/go-pow/value-types"
)

type ChallengeHashData struct {
	Challenge      Challenge
	TargetBitIndex int
	Nonce          powValueTypes.Nonce
}

type Challenge struct {
	leadingZeroCount powValueTypes.LeadingZeroCount
	createdAt        mo.Option[powValueTypes.CreatedAt]
	resource         mo.Option[powValueTypes.Resource]
	payload          powValueTypes.Payload
	hash             powValueTypes.Hash
	hashDataLayout   powValueTypes.HashDataLayout
}

func (entity Challenge) LeadingZeroCount() powValueTypes.LeadingZeroCount {
	return entity.leadingZeroCount
}

func (entity Challenge) TargetBitIndex() int {
	return entity.hash.SizeInBits() - entity.leadingZeroCount.ToInt()
}

func (entity Challenge) CreatedAt() mo.Option[powValueTypes.CreatedAt] {
	return entity.createdAt
}

func (entity Challenge) Resource() mo.Option[powValueTypes.Resource] {
	return entity.resource
}

func (entity Challenge) Payload() powValueTypes.Payload {
	return entity.payload
}

func (entity Challenge) Hash() powValueTypes.Hash {
	return entity.hash
}

func (entity Challenge) HashDataLayout() powValueTypes.HashDataLayout {
	return entity.hashDataLayout
}

func (entity Challenge) Solve() (Solution, error) {
	targetBitIndex := entity.TargetBitIndex()
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
