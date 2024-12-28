package utils

import (
	"math/big"
	"sort"
)

func CalcNakamotoCoefficientBigNums(totalVotingPower *big.Int, votingPowers []big.Int) int {
	thresholdPercent := big.NewFloat(0.33)
	thresholdVal := new(big.Float).Mul(new(big.Float).SetInt(totalVotingPower), thresholdPercent)
	cumulativeVal := big.NewFloat(0.00)
	nakamotoCoefficient := 0

	sort.Slice(votingPowers, func(i, j int) bool {
		return votingPowers[i].Cmp(&votingPowers[j]) > 0
	})

	for _, vp := range votingPowers {
		z := new(big.Float).Add(cumulativeVal, new(big.Float).SetInt(&vp))
		cumulativeVal = z
		nakamotoCoefficient += 1
		if cumulativeVal.Cmp(thresholdVal) == +1 {
			break
		}
	}

	return nakamotoCoefficient
}

// CalcNakamotoCoefficientBigNums51 calculates the Nakamoto Coefficient for a given voting power distribution.
// This function is used for chains that require at least > 50% of the total voting power to be controlled,
// that therefore cannot use the 33% threshold with the CalcNakamotoCoefficientBigNums function. 
func CalcNakamotoCoefficientBigNums51(totalVotingPower *big.Int, votingPowers []big.Int) int {
	thresholdPercent := big.NewFloat(0.50)
	thresholdVal := new(big.Float).Mul(new(big.Float).SetInt(totalVotingPower), thresholdPercent)
	cumulativeVal := big.NewFloat(0.00)
	nakamotoCoefficient := 0

	for _, vp := range votingPowers {
		z := new(big.Float).Add(cumulativeVal, new(big.Float).SetInt(&vp))
		cumulativeVal = z
		nakamotoCoefficient += 1
		if cumulativeVal.Cmp(thresholdVal) == +1 {
			break
		}
	}

	return nakamotoCoefficient
}