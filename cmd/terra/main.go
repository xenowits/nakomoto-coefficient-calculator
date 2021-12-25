package terra

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	utils "github.com/xenowits/nakamoto-coefficient-calculator/utils"
)

type Request struct {
	height   int
	page     int
	per_page int
}

type Response struct {
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
		Total string `json:"total"`
	} `json:"result"`
}

type ErrorResponse struct {
	Id      int    `json:"id"`
	Jsonrpc string `json:"jsonrpc"`
	Error   string `json:"error"`
}

// https://fcd.terra.dev/swagger
func Terra() (int, error) {
	votingPowers := make([]int64, 0, 200)
	url := fmt.Sprintf("https://fcd.terra.dev/validatorsets/latest")
	resp, err := http.Get(url)
	if err != nil {
		errBody, _ := ioutil.ReadAll(resp.Body)
		var errResp ErrorResponse
		json.Unmarshal(errBody, &errResp)
		log.Println(errResp.Error)
		return -1, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return -1, err
	}

	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		return -1, err
	}

	// loop through the validators voting powers
	for _, ele := range response.Result.Validators {
		val, _ := strconv.Atoi(ele.Voting_power)
		votingPowers = append(votingPowers, int64(val))
	}

	// No need to sort as the result is already in sorted in descending order
	totalVotingPower := utils.CalculateTotalVotingPower(votingPowers)
	fmt.Println("Total voting power:", totalVotingPower)

	// now we're ready to calculate the nakomoto coefficient
	nakamotoCoefficient := utils.CalcNakamotoCoefficient(totalVotingPower, votingPowers)
	fmt.Println("The Nakamoto coefficient for terra is", nakamotoCoefficient)

	return nakamotoCoefficient, nil
}
