package utils

import (
	"math/big"
)

// THRESHOLD for calculating the Nakamoto coefficient (33%)
const THRESHOLD = 3

// CalculateTotalVotingPowerBigInt calculates the total voting power from a slice of big.Int
func CalculateTotalVotingPowerBigInt(votingPowers []*big.Int) *big.Int {
	total := big.NewInt(0)
	for _, power := range votingPowers {
		total.Add(total, power)
	}
	return total
}

// CalcNakamotoCoefficientBigInt calculates the Nakamoto coefficient using big.Int
func CalcNakamotoCoefficientBigInt(totalVotingPower *big.Int, votingPowers []*big.Int) int {
	threshold := new(big.Int).Div(totalVotingPower, big.NewInt(THRESHOLD)) // 33% threshold
	cumulativePower := big.NewInt(0)
	for i, power := range votingPowers {
		cumulativePower.Add(cumulativePower, power)
		if cumulativePower.Cmp(threshold) >= 0 {
			return i + 1
		}
	}
	return len(votingPowers)
}
