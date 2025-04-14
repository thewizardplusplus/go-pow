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
	ttl                 mo.Option[powValueTypes.TTL]
	resource            mo.Option[powValueTypes.Resource]
	serializedPayload   mo.Option[powValueTypes.SerializedPayload]
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

func (builder *ChallengeBuilder) SetTTL(
	value powValueTypes.TTL,
) *ChallengeBuilder {
	builder.ttl = mo.Some(value)
	return builder
}

func (builder *ChallengeBuilder) SetResource(
	value powValueTypes.Resource,
) *ChallengeBuilder {
	builder.resource = mo.Some(value)
	return builder
}

func (builder *ChallengeBuilder) SetSerializedPayload(
	value powValueTypes.SerializedPayload,
) *ChallengeBuilder {
	builder.serializedPayload = mo.Some(value)
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

	if builder.createdAt.IsPresent() != builder.ttl.IsPresent() {
		errs = append(
			errs,
			errors.New(
				"`CreatedAt` timestamp and TTL "+
					"should either both be specified or both omitted",
			),
		)
	}

	serializedPayload, isPresent := builder.serializedPayload.Get()
	if !isPresent {
		errs = append(errs, errors.New("serialized payload is required"))
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
		ttl:                 builder.ttl,
		resource:            builder.resource,
		serializedPayload:   serializedPayload,
		hash:                hash,
		hashDataLayout:      hashDataLayout,
	}
	if err := builder.checkHashDataLayout(entity); err != nil {
		return Challenge{}, fmt.Errorf(
			"unable to check the hash data layout: %w",
			err,
		)
	}

	return entity, nil
}

func (builder ChallengeBuilder) checkHashDataLayout(entity Challenge) error {
	nonce, err := powValueTypes.NewZeroNonce()
	if err != nil {
		return fmt.Errorf("unable to construct the zero nonce: %w", err)
	}

	if _, err := entity.hashDataLayout.Execute(ChallengeHashData{
		Challenge: entity,
		Nonce:     nonce,
	}); err != nil {
		return fmt.Errorf("unable to execute the hash data layout: %w", err)
	}

	return nil
}
