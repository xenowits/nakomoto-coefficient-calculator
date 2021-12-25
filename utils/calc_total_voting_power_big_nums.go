package utils

import (
	"math/big"
)
func CalculateTotalVotingPowerBigNums(votingPowers []big.Int) *big.Int {
	total := big.NewInt(0)
	for _, vp := range votingPowers {
		total = new(big.Int).Add(total, &vp)
	}
	return total
}