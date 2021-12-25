package polygon

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"

	utils "github.com/xenowits/nakamoto-coefficient-calculator/utils"
)

type Response struct {
	Height string `json:"height"`
	Result struct {
		Block_height string
		Validators   []struct {
			ID         int    `json:"id"`
			StartEpoch int    `json:"startEpoch"`
			EndEpoch   int    `json:"endEpoch"`
			Nonce      int    `json:"nonce"`
			Power      int64  `json:"power"`
			PubKey     string `json:"pubKey"`
			Signer     string `json:"signer"`
			Jailed     bool   `json:"jailed"`
		} `json:"validators"`
		Count string `json:"count"`
		Total string `json:"total"`
	} `json:"result"`
}

type ErrorResponse struct {
	Id      int    `json:"id"`
	Jsonrpc string `json:"jsonrpc"`
	Error   string `json:"error"`
}

func Polygon() (int, error) {
	votingPowers := make([]int64, 0, 200)

	url := fmt.Sprintf("https://heimdall.api.matic.network/staking/validator-set")
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
		votingPowers = append(votingPowers, int64(ele.Power))
	}

	// need to sort the powers in descending order since they are in random order
	sort.Slice(votingPowers, func(i, j int) bool { return votingPowers[i] > votingPowers[j] })

	totalVotingPower := utils.CalculateTotalVotingPower(votingPowers)
	fmt.Println("Total voting power:", totalVotingPower)

	// // now we're ready to calculate the nakomoto coefficient
	nakamotoCoefficient := utils.CalcNakamotoCoefficient(totalVotingPower, votingPowers)
	fmt.Println("The Nakamoto coefficient for 0xPolygon is", nakamotoCoefficient)

	return nakamotoCoefficient, nil
}
