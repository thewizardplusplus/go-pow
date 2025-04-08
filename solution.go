package pow

import (
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
