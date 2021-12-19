package utils

func CalculateTotalVotingPower(votingPowers []int64) int64 {
	var total int64 = 0
	for _, vp := range votingPowers {
		total += vp
	}
	return total
}