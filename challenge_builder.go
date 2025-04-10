package pow

import (
	"errors"
	"fmt"

	"github.com/samber/mo"
	powValueTypes "github.com/thewizardplusplus/go-pow/value-types"
)

type ChallengeBuilder struct {
	leadingZeroBitCount mo.Option[powValueTypes.LeadingZeroBitCount]
	targetBitIndex      mo.Option[powValueTypes.TargetBitIndex]
	createdAt           mo.Option[powValueTypes.CreatedAt]
	resource            mo.Option[powValueTypes.Resource]
	payload             mo.Option[powValueTypes.Payload]
	hash                mo.Option[powValueTypes.Hash]
	hashDataLayout      mo.Option[powValueTypes.HashDataLayout]
}

func NewChallengeBuilder() *ChallengeBuilder {
	return &ChallengeBuilder{}
}

func (builder *ChallengeBuilder) SetLeadingZeroBitCount(
	value powValueTypes.LeadingZeroBitCount,
) *ChallengeBuilder {
	builder.leadingZeroBitCount = mo.Some(value)
	return builder
}

func (builder *ChallengeBuilder) SetTargetBitIndex(
	value powValueTypes.TargetBitIndex,
) *ChallengeBuilder {
	builder.targetBitIndex = mo.Some(value)
	return builder
}

func (builder *ChallengeBuilder) SetCreatedAt(
	value powValueTypes.CreatedAt,
) *ChallengeBuilder {
	builder.createdAt = mo.Some(value)
	return builder
}

func (builder *ChallengeBuilder) SetResource(
	value powValueTypes.Resource,
) *ChallengeBuilder {
	builder.resource = mo.Some(value)
	return builder
}

func (builder *ChallengeBuilder) SetPayload(
	value powValueTypes.Payload,
) *ChallengeBuilder {
	builder.payload = mo.Some(value)
	return builder
}

func (builder *ChallengeBuilder) SetHash(
	value powValueTypes.Hash,
) *ChallengeBuilder {
	builder.hash = mo.Some(value)
	return builder
}

func (builder *ChallengeBuilder) SetHashDataLayout(
	value powValueTypes.HashDataLayout,
) *ChallengeBuilder {
	builder.hashDataLayout = mo.Some(value)
	return builder
}

func (builder ChallengeBuilder) Build() (Challenge, error) {
	var errs []error

	leadingZeroBitCount, isLeadingZeroBitCountPresent :=
		builder.leadingZeroBitCount.Get()
	targetBitIndex, isTargetBitIndexPresent := builder.targetBitIndex.Get()
	if !isLeadingZeroBitCountPresent && !isTargetBitIndexPresent {
		errs = append(
			errs,
			errors.New("leading zero bit count or target bit index is required"),
		)
	} else if isLeadingZeroBitCountPresent && isTargetBitIndexPresent {
		errs = append(
			errs,
			errors.New(
				"leading zero bit count and target bit index "+
					"are specified at the same time",
			),
		)
	}

	payload, isPresent := builder.payload.Get()
	if !isPresent {
		errs = append(errs, errors.New("payload is required"))
	}

	hash, isPresent := builder.hash.Get()
	if !isPresent {
		errs = append(errs, errors.New("hash is required"))
	} else if isLeadingZeroBitCountPresent &&
		leadingZeroBitCount.ToInt() > hash.SizeInBits() {
		errs = append(
			errs,
			errors.New("leading zero bit count exceeds the hash checksum size"),
		)
	} else if isTargetBitIndexPresent {
		rawLeadingZeroBitCount := hash.SizeInBits() - targetBitIndex.ToInt()

		var err error
		leadingZeroBitCount, err = powValueTypes.NewLeadingZeroBitCount(
			rawLeadingZeroBitCount,
		)
		if err != nil {
			errs = append(
				errs,
				fmt.Errorf("unable to construct the leading zero bit count: %w", err),
			)
		}
	}

	hashDataLayout, isPresent := builder.hashDataLayout.Get()
	if !isPresent {
		errs = append(errs, errors.New("hash data layout is required"))
	}

	if len(errs) > 0 {
		return Challenge{}, errors.Join(errs...)
	}

	entity := Challenge{
		leadingZeroBitCount: leadingZeroBitCount,
		createdAt:           builder.createdAt,
		resource:            builder.resource,
		payload:             payload,
		hash:                hash,
		hashDataLayout:      hashDataLayout,
	}
	return entity, nil
}
