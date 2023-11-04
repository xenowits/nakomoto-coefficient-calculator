package chains

import (
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"sort"
	"strconv"

	utils "github.com/xenowits/nakamoto-coefficient-calculator/core/utils"
)

const minVotingBalance int64 = 32000000000000000

type Validator struct {
	Balance string `json:"balance"`
	Status  string `json:"status"`
}

type PulsechainResponse struct {
	Validators []Validator `json:"data"`
}

type PulsechainErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
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
	var votingPowers []big.Int

	url := fmt.Sprintf("https://rpc-pulsechain.g4mm4.io/beacon-api/eth/v1/beacon/states/head/validators")
	resp, err := http.Get(url)
	if err != nil {
		errBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return 0, err
		}
		var errResp PulsechainErrorResponse

		errr := json.Unmarshal(errBody, &errResp)
		if errr != nil {
			return 0, errr
		}

		return 0, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
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
		balance, err := strconv.ParseInt(validator.Balance, 10, 64)

		return err == nil && validator.Status == "active_ongoing" &&
			balance >= minVotingBalance
	})

	// loop through the validators again, create the votingPowers array
	for _, validator := range activeValidators {
		balance, err := strconv.ParseInt(validator.Balance, 10, 64)
		if err != nil {
			return 0, err
		}

		votingPowers = append(votingPowers, *big.NewInt(balance))
	}

	// Sort the voting powers in descending order since they maybe in random order.
	// Sort the powers in descending order since they maybe in random order
	sort.Slice(votingPowers, func(i, j int) bool {
		res := (&votingPowers[i]).Cmp(&votingPowers[j])
		return res == 1
	})

	totalVotingPower := utils.CalculateTotalVotingPowerBigNums(votingPowers)
	fmt.Println("Total voting power:", totalVotingPower)

	// Now we're ready to calculate the nakamoto coefficient
	nakamotoCoefficient := utils.CalcNakamotoCoefficientBigNums(totalVotingPower, votingPowers)
	fmt.Println("The Nakamoto coefficient for Pulsechain is", nakamotoCoefficient)

	return nakamotoCoefficient, nil
}
