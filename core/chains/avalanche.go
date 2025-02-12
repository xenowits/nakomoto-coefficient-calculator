package chains

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"sort"
	_ "strings"
)

type AvalancheResponse struct {
	Jsonrpc string `json:"jsonrpc"`
	Id      int    `json:"id"`
	Result  struct {
		Validators []struct {
			Weight string `json:"weight"` // Correct field for stake amount
		} `json:"validators"`
	} `json:"result"`
}

// Avalanche calculates the Nakamoto coefficient for Avalanche C-Chain.
func Avalanche() (int, error) {
	var votingPowers []*big.Int

	url := "https://api.avax.network/ext/P"
	jsonReqData := []byte(`{"jsonrpc": "2.0","method": "platform.getCurrentValidators","params":{},"id":1}`)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonReqData))
	if err != nil {
		return 0, fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return 0, fmt.Errorf("API request failed with status code %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("failed to read response body: %v", err)
	}

	var response AvalancheResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return 0, fmt.Errorf("failed to parse JSON response: %v", err)
	}

	if len(response.Result.Validators) == 0 {
		return 0, fmt.Errorf("no validators found in API response")
	}

	// Parse stake amounts from "weight" field and compute total voting power
	totalVotingPower := big.NewInt(0)
	for _, v := range response.Result.Validators {
		if v.Weight == "" {
			continue
		}

		stake := new(big.Int)
		stakeFloat := new(big.Float)

		if _, success := stakeFloat.SetString(v.Weight); success {
			stakeFloat.Int(stake)
		} else if _, success := stake.SetString(v.Weight, 10); !success {
			continue
		}

		votingPowers = append(votingPowers, stake)
		totalVotingPower.Add(totalVotingPower, stake)
	}

	if totalVotingPower.Cmp(big.NewInt(0)) == 0 {
		return 0, fmt.Errorf("total voting power is still 0, check API response")
	}

	// Sort voting powers in descending order
	sort.Slice(votingPowers, func(i, j int) bool {
		return votingPowers[i].Cmp(votingPowers[j]) > 0
	})

	// Calculate Nakamoto coefficient (top n validators > 33% of stake)
	threshold := new(big.Int).Div(new(big.Int).Mul(totalVotingPower, big.NewInt(33)), big.NewInt(100)) // 33% of total stake
	accumulatedPower := big.NewInt(0)
	nakamotoCoefficient := 0

	for _, power := range votingPowers {
		accumulatedPower.Add(accumulatedPower, power)
		nakamotoCoefficient++

		if accumulatedPower.Cmp(threshold) >= 0 {
			break
		}
	}

	// Final result output
	fmt.Println("Total voting power:", totalVotingPower)
	fmt.Println("The Nakamoto coefficient for Avalanche is", nakamotoCoefficient)

	return nakamotoCoefficient, nil
}
