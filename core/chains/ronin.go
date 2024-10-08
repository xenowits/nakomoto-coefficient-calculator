package chains

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/big"
	"net/http"
	"sort"
	"strings"

	"github.com/xenowits/nakamoto-coefficient-calculator/core/utils"
)

type RoninResponse []struct {
	Address     string `json:"address"`
	TotalStaked string `json:"totalStaked"`
	Status      string `json:"status"`
}

const MaxFinalityVotePercentage uint16 = 10_000
const MaxFinalityVoteThreshold = 22

const RoninValidatorQuery = `{
  "query": "query ValidatorOrCandidates { ValidatorOrCandidates { address totalStaked status } }"
}`

func Ronin() (int, error) {

	url := fmt.Sprintf("https://indexer.roninchain.com/query")

	var totalStakes []*big.Int

	// Create a new request using http
	req, err := http.NewRequest("POST", url, strings.NewReader(RoninValidatorQuery))

	// Send req using http Client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	var response struct {
		Data struct {
			ValidatorOrCandidates RoninResponse
		}
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return 0, err
	}

	// loop through the validators stakes
	for _, ele := range response.Data.ValidatorOrCandidates {
		// skip validators that are no longer active
		if ele.Status != "" {
			continue
		}
		staked := new(big.Int)
		staked, ok := staked.SetString(ele.TotalStaked, 10)
		if !ok {
			return 0, fmt.Errorf("failed to convert %s to big.Int", ele.TotalStaked)
		}
		totalStakes = append(totalStakes, staked)
	}

	// need to sort the stakes in descending order since they are in random order
	sort.Slice(totalStakes, func(i, j int) bool {
		res := (totalStakes[i]).Cmp(totalStakes[j])
		return res == 1
	})

	// normalizes the finality vote weight
	_ = normalizeFinalityVoteWeight(totalStakes, MaxFinalityVoteThreshold)

	// convert []*big.Int to []big.Int
	totalStakesResult := pointerSliceToSlice(totalStakes)

	totalStake := utils.CalculateTotalVotingPowerBigNums(totalStakesResult)
	fmt.Println("Total voting power:", new(big.Float).SetInt(totalStake))

	// now we're ready to calculate the nakomoto coefficient
	nakamotoCoefficient := utils.CalcNakamotoCoefficientBigNums(totalStake, totalStakesResult)
	fmt.Println("The Nakamoto coefficient for Ronin is", nakamotoCoefficient)

	return nakamotoCoefficient, nil
}

func pointerSliceToSlice(s []*big.Int) []big.Int {
	var slice []big.Int
	for _, v := range s {
		slice = append(slice, *v)
	}
	return slice
}

// From: https://github.com/axieinfinity/ronin/blob/adc5849b5af046532755e5dae8f023c143d5494c/consensus/consortium/common/utils.go#L76
func normalizeFinalityVoteWeight(stakedAmounts []*big.Int, threshold int) []uint16 {
	weights := make([]uint16, 0, len(stakedAmounts))

	// The candidate list is too small, so weight is equal among the candidates
	if len(stakedAmounts) <= threshold {
		for range stakedAmounts {
			weights = append(weights, MaxFinalityVotePercentage/uint16(len(stakedAmounts)))
		}

		return weights
	}

	cpyStakedAmounts := make([]*big.Int, len(stakedAmounts))
	for i, stakedAmount := range stakedAmounts {
		cpyStakedAmounts[i] = new(big.Int).Set(stakedAmount)
	}

	// Sort staked amount in descending order
	for i := 0; i < len(cpyStakedAmounts)-1; i++ {
		for j := i + 1; j < len(cpyStakedAmounts); j++ {
			if cpyStakedAmounts[i].Cmp(cpyStakedAmounts[j]) < 0 {
				cpyStakedAmounts[i], cpyStakedAmounts[j] = cpyStakedAmounts[j], cpyStakedAmounts[i]
			}
		}
	}

	totalStakedAmount := new(big.Int)
	for _, stakedAmount := range cpyStakedAmounts {
		totalStakedAmount.Add(totalStakedAmount, stakedAmount)
	}
	weightThreshold := new(big.Int).Div(totalStakedAmount, big.NewInt(int64(threshold)))

	pointer := 0
	sumOfUnchangedElements := totalStakedAmount
	for {
		sumOfChangedElements := new(big.Int)
		shouldBreak := true
		for cpyStakedAmounts[pointer].Cmp(weightThreshold) > 0 {
			sumOfChangedElements.Add(sumOfChangedElements, cpyStakedAmounts[pointer])
			shouldBreak = false
			pointer++
		}

		if shouldBreak {
			break
		}

		sumOfUnchangedElements = new(big.Int).Sub(sumOfUnchangedElements, sumOfChangedElements)
		weightThreshold = new(big.Int).Div(
			sumOfUnchangedElements,
			new(big.Int).Sub(big.NewInt(int64(threshold)), big.NewInt(int64(pointer))),
		)
	}

	for i, stakedAmount := range stakedAmounts {
		if stakedAmount.Cmp(weightThreshold) > 0 {
			stakedAmounts[i] = weightThreshold
		}
	}

	totalStakedAmount.SetUint64(0)
	for _, stakedAmount := range stakedAmounts {
		totalStakedAmount.Add(totalStakedAmount, stakedAmount)
	}

	for _, stakedAmount := range stakedAmounts {
		weight := new(big.Int).Mul(stakedAmount, big.NewInt(int64(MaxFinalityVotePercentage)))
		weight.Div(weight, totalStakedAmount)

		weights = append(weights, uint16(weight.Uint64()))
	}

	// Due to the imprecision of division, the remaining weight for the total to reach 100% is
	// split equally across cnadidates. After this step, the total weight may still not reach
	// 100% but the imprecision is neglectible (lower than the length of candidate list)
	var totalFinalityWeight uint16
	for _, weight := range weights {
		totalFinalityWeight += weight
	}
	cutOffWeight := MaxFinalityVotePercentage - totalFinalityWeight
	topUpWeight := cutOffWeight / uint16(len(weights))
	for i := range weights {
		weights[i] += topUpWeight
	}

	return weights
}
