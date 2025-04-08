package pow

import (
	"errors"

	"github.com/samber/mo"
	powValueTypes "github.com/thewizardplusplus/go-pow/value-types"
)

type SolutionBuilder struct {
	challenge mo.Option[Challenge]
	nonce     mo.Option[powValueTypes.Nonce]
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

func (builder SolutionBuilder) Build() (Solution, error) {
	var errs []error

	challenge, isPresent := builder.challenge.Get()
	if !isPresent {
		errs = append(errs, errors.New("challenge is required"))
	}

	nonce, isPresent := builder.nonce.Get()
	if !isPresent {
		errs = append(errs, errors.New("nonce is required"))
	}

	if len(errs) > 0 {
		return Solution{}, errors.Join(errs...)
	}

	entity := Solution{
		challenge: challenge,
		nonce:     nonce,
	}
	return entity, nil
}
