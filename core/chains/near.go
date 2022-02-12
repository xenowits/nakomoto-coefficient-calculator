package chains

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/xenowits/nakamoto-coefficient-calculator/core/utils"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"sort"
)

type NearResponse struct {
	Jsonrpc string `json:"jsonrpc"`
	Id      int    `json:"id"`
	Result  struct {
		Validators []struct {
			AccountId string `json:"account_id"`
			Stake     string `json:"stake"`
		} `json:"current_validators"`
	} `json:"result"`
}

func Near() (int, error) {
	votingPowers := make([]big.Int, 0, 1024)

	url := fmt.Sprintf("https://rpc.mainnet.near.org")
	jsonReqData := []byte(`{"jsonrpc": "2.0","method": "validators","params":[null],"id":1}`)

	// Create a new POST request using http
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonReqData))
	if err != nil {
		return -1, err
	}
	req.Header.Set("Content-Type", "application/json")

	// Send req using http Client
	client := &http.Client{}
	resp, err := client.Do(req)

	if resp != nil {
		// Need to close body when redirection occurs
		// In redirection, response is not empty
		defer func() {
			if err := resp.Body.Close(); err != nil {
				log.Printf("failed to close response body")
			}
		}()
	}

	if err != nil {
		return -1, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return -1, err
	}

	var response NearResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return -1, err
	}

	// loop through the validators voting powers
	for _, ele := range response.Result.Validators {
		n, ok := new(big.Int).SetString(ele.Stake, 10)
		if !ok {
			return -1, fmt.Errorf("failed to parse string %s", ele.Stake)
		}
		votingPowers = append(votingPowers, *n)
	}

	// need to sort the powers in descending order since they are in random order
	sort.Slice(votingPowers, func(i, j int) bool {
		res := (&votingPowers[i]).Cmp(&votingPowers[j])
		if res == 1 {
			return true
		}
		return false
	})

	totalVotingPower := utils.CalculateTotalVotingPowerBigNums(votingPowers)
	fmt.Println("Total voting power:", totalVotingPower)

	// now we're ready to calculate the nakamoto coefficient
	nakamotoCoefficient := utils.CalcNakamotoCoefficientBigNums(totalVotingPower, votingPowers)
	fmt.Println("The Nakamoto coefficient for near protocol is", nakamotoCoefficient)

	return nakamotoCoefficient, nil
}
