package pow

import (
	"errors"

	"github.com/samber/mo"
	powValueTypes "github.com/thewizardplusplus/go-pow/value-types"
)

type SolutionBuilder struct {
	challenge mo.Option[Challenge]
	nonce     mo.Option[powValueTypes.Nonce]
	hashSum   mo.Option[powValueTypes.HashSum]
}

func NewSolutionBuilder() *SolutionBuilder {
	return &SolutionBuilder{}
}

func (builder *SolutionBuilder) SetChallenge(value Challenge) *SolutionBuilder {
	builder.challenge = mo.Some(value)
	return builder
}

func (builder *SolutionBuilder) SetNonce(
	value powValueTypes.Nonce,
) *SolutionBuilder {
	builder.nonce = mo.Some(value)
	return builder
}

func (builder *SolutionBuilder) SetHashSum(
	value powValueTypes.HashSum,
) *SolutionBuilder {
	builder.hashSum = mo.Some(value)
	return builder
}

func (builder SolutionBuilder) Build() (Solution, error) {
	var errs []error

	challenge, isChallengePresent := builder.challenge.Get()
	if !isChallengePresent {
		errs = append(errs, errors.New("challenge is required"))
	}

	nonce, isPresent := builder.nonce.Get()
	if !isPresent {
		errs = append(errs, errors.New("nonce is required"))
	}

	hashSum, isHashSumPresent := builder.hashSum.Get()
	if isHashSumPresent &&
		isChallengePresent &&
		hashSum.Len() != challenge.hash.SizeInBytes() {
		errs = append(
			errs,
			errors.New("hash sum length doesn't match the hash checksum size"),
		)
	}

	if len(errs) > 0 {
		return Solution{}, errors.Join(errs...)
	}

	entity := Solution{
		challenge: challenge,
		nonce:     nonce,
		hashSum:   builder.hashSum,
	}
	return entity, nil
}
