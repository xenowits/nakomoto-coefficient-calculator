package chains

import (
	"encoding/json"
	"fmt"
	"strconv"
	"io/ioutil"
	"log"
	"sort"
	"net/http"

	utils "github.com/xenowits/nakamoto-coefficient-calculator/core/utils"
)

type Validator struct {
	Balance   string     `json:"balance"`
	Status    string    `json:"status"`
}

type PulsechainResponse struct {
	Validators []Validator `json:"data"`
}

type PulsechainErrorResponse struct {
	Code      int      `json:"code"`
	Message   string   `json:"message"`
}

func Filter[T any](s []T, cond func(t T) bool) []T {
	res := []T{}
	for _, v := range s {
		if cond(v) {
			res = append(res, v)
		}
	}
	return res
}

func Pulsechain() (int, error) {
	var votingPowers []int64
	
	url := fmt.Sprintf("https://rpc-pulsechain.g4mm4.io/beacon-api/eth/v1/beacon/states/head/validators")
	resp, err := http.Get(url)
	if err != nil {
		errBody, _ := ioutil.ReadAll(resp.Body)
		var errResp PulsechainErrorResponse
		json.Unmarshal(errBody, &errResp)
		return 0, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	var response PulsechainResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return 0, err
	}

	// break if no entries
	if len(response.Validators) == 0 {
		return 0, err
	}

	activeValidators := Filter(response.Validators, func(validator Validator) bool {
		return validator.Status == "active_ongoing"
	})

	// loop through the validators again, create the votingPowers array
	for _, validator := range activeValidators {
		balance, _ := strconv.ParseInt(validator.Balance, 10, 64)
		votingPowers = append(votingPowers, balance)
	}

	// Sort the voting powers in descending order since they maybe in random order.
	sort.Slice(votingPowers, func(i, j int) bool { return votingPowers[i] > votingPowers[j] })

	totalVotingPower := utils.CalculateTotalVotingPower(votingPowers)
	fmt.Println("Total voting power:", totalVotingPower)

	// Now we're ready to calculate the nakamoto coefficient
	nakamotoCoefficient := utils.CalcNakamotoCoefficient(totalVotingPower, votingPowers)
	fmt.Println("The Nakamoto coefficient for Pulsechain is", nakamotoCoefficient)

	return nakamotoCoefficient, nil
}
