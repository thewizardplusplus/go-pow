package pow

import (
	"github.com/samber/mo"
	powValueTypes "github.com/thewizardplusplus/go-pow/value-types"
)

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
