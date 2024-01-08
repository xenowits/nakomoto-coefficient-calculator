package utils

import "sort"

func CalcNakamotoCoefficient(totalVotingPower int64, votingPowers []int64) int {
	var cumulativePercent, thresholdPercent, curr float64 = 0.00, THRESHOLD_PERCENT, 0.00
	nakamotoCoefficient := 0

	sort.Slice(votingPowers, func(i, j int) bool { return votingPowers[i] > votingPowers[j] })

	for _, vp := range votingPowers {
		curr = float64(vp) / float64(totalVotingPower)
		cumulativePercent += curr * 100
		nakamotoCoefficient += 1
		if cumulativePercent >= thresholdPercent {
			break
		}
	}
	return nakamotoCoefficient
}
