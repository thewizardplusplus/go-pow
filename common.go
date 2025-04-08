package pow

import (
	"math/big"

	powValueTypes "github.com/thewizardplusplus/go-pow/value-types"
)

func makeTarget(targetBitIndex int) *big.Int {
	target := big.NewInt(0)
	target.SetBit(target, targetBitIndex, 1)

	return target
}

func isHashSumFitTarget(hashSum powValueTypes.HashSum, target *big.Int) bool {
	hashSumAsBigInt := big.NewInt(0)
	hashSumAsBigInt.SetBytes(hashSum.ToBytes())

	return hashSumAsBigInt.Cmp(target) == -1 // hashSumAsBigInt < target
}
