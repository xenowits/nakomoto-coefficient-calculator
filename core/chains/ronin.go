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

const RoninValidatorQuery = `{
  "query": "query ValidatorOrCandidates { ValidatorOrCandidates { address totalStaked status } }"
}`

func Ronin() (int, error) {

	url := fmt.Sprintf("https://indexer.roninchain.com/query")

	var totalStakes []big.Int

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

	// loop through the validators voting powers
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
		totalStakes = append(totalStakes, *staked)
	}

	// need to sort the powers in descending order since they are in random order
	sort.Slice(totalStakes, func(i, j int) bool {
		res := (&totalStakes[i]).Cmp(&totalStakes[j])
		if res == 1 {
			return true
		}
		return false
	})

	totalStake := utils.CalculateTotalVotingPowerBigNums(totalStakes)
	fmt.Println("Total voting power:", new(big.Float).SetInt(totalStake))

	// now we're ready to calculate the nakomoto coefficient
	nakamotoCoefficient := utils.CalcNakamotoCoefficientBigNums(totalStake, totalStakes)
	fmt.Println("The Nakamoto coefficient for Ronin is", nakamotoCoefficient)

	return nakamotoCoefficient, nil
}
