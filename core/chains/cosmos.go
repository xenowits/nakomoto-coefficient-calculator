package chains

import (
	"encoding/json"
	"fmt"
	utils "github.com/xenowits/nakamoto-coefficient-calculator/core/utils"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type CosmosResponse struct {
	Jsonrpc string `json:"jsonrpc"`
	Id      int    `json:"id"`
	Result struct {
		Block_height string
		Validators   []struct {
			Address           string `json:"address"`
			Voting_power      string `json:"voting_power"`
			Proposer_priority string `json:"proposer_priority"`
			Pub_key           struct {
				Type  string `json:"type"`
				Value string `json:"value"`
			} `json:"pub_key"`
		} `json:"validators"`
		Count string `json:"count"`
		Total string `json:"total"`
	} `json:"result"`
}

type CosmosErrorResponse struct {
	Id      int    `json:"id"`
	Jsonrpc string `json:"jsonrpc"`
	Error   string `json:"error"`
}

func Cosmos() (int, error) {
	votingPowers := make([]int64, 0, 200)
	pageNo, entriesPerPage := 1, 50
	url := ""
	for true {
		url = fmt.Sprintf("https://rpc.cosmos.network/validators?page=%d&per_page=%d", pageNo, entriesPerPage)
		resp, err := http.Get(url)
		if err != nil {
			errBody, _ := ioutil.ReadAll(resp.Body)
			var errResp CosmosErrorResponse
			json.Unmarshal(errBody, &errResp)
			log.Println(errResp.Error)
			return -1, nil
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return -1, nil
		}

		var response CosmosResponse
		err = json.Unmarshal(body, &response)
		if err != nil {
			return -1, nil
		}

		// break if no more entries left
		if len(response.Result.Validators) == 0 {
			break
		}

		// loop through the validators voting powers
		for _, ele := range response.Result.Validators {
			val, _ := strconv.Atoi(ele.Voting_power)
			votingPowers = append(votingPowers, int64(val))
		}

		// increment counters
		pageNo += 1
	}

	// No need to sort as the result is already in sorted in descending order
	totalVotingPower := utils.CalculateTotalVotingPower(votingPowers)
	fmt.Println("Total voting power:", totalVotingPower)

	// now we're ready to calculate the nakomoto coefficient
	nakamotoCoefficient := utils.CalcNakamotoCoefficient(totalVotingPower, votingPowers)
	fmt.Println("The Nakamoto coefficient for cosmos is", nakamotoCoefficient)

	return nakamotoCoefficient, nil
}
