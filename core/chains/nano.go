package chains

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	utils "github.com/xenowits/nakamoto-coefficient-calculator/core/utils"
	"io/ioutil"
	"math/big"
	"net/http"
	"sort"
)

type NanoResponse struct {
	Representatives map[string]string `json:"representatives"`
}

func Nano() (int, error) {
	votingPowers := make([]big.Int, 0, 1000)

	// node docs at https://powernode.cc/api
	// protocol docs at https://docs.nano.org/commands/rpc-protocol/
	url := fmt.Sprintf("https://proxy.powernode.cc/proxy")
	reqData := map[string]string{"action": "representatives", "sorting": "true"}
	reqJson, err := json.Marshal(reqData)
	if err != nil {
		return -1, err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(reqJson))
	if err != nil {
		return -1, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return -1, err
	}

	var response NanoResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return -1, err
	}
	
	zero := big.NewInt(0)
	// loop through the validators voting powers
	for _, weight := range response.Representatives {
		n := new(big.Int)
		n, ok := n.SetString(weight, 10)
		if !ok {
			return -1, errors.New("invalid number string")
		}
		if n.Cmp(zero) == 1 {
			// ignore zero balance validators
			// shouldn't change the result, but saves a tiny bit of overhead
			votingPowers = append(votingPowers, *n)
		}
	}

	// sort the powers in descending order
	sort.Slice(votingPowers, func(i, j int) bool {
		res := (&votingPowers[i]).Cmp(&votingPowers[j])
		if res == 1 {
			return true
		}
		return false
	})

	totalVotingPower := utils.CalculateTotalVotingPowerBigNums(votingPowers)
	fmt.Println("Total voting power:", new(big.Float).SetInt(totalVotingPower))

	// calculate the nakomoto coefficient
	nakamotoCoefficient := utils.CalcNakamotoCoefficientBigNums(totalVotingPower, votingPowers)
	fmt.Println("The Nakamoto coefficient for Nano is", nakamotoCoefficient)

	return nakamotoCoefficient, nil
}
