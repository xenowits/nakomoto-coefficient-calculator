package chains

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"
	"strconv"

	utils "github.com/xenowits/nakamoto-coefficient-calculator/core/utils"
)

type AgoricResponse struct {
	Result []struct {
		Tokens      string `json:"tokens"`
	} `json:"result"`
}

type AgoricErrorResponse struct {
	Id      int    `json:"id"`
	Jsonrpc string `json:"jsonrpc"`
	Error   string `json:"error"`
}

func Agoric() (int, error) {
	votingPowers := make([]int64, 0, 200)

	url := fmt.Sprintf("https://lcd-agoric.keplr.app/staking/validators")
	resp, err := http.Get(url)
	if err != nil {
		errBody, _ := io.ReadAll(resp.Body)
		var errResp AgoricErrorResponse
		json.Unmarshal(errBody, &errResp)
		log.Println(errResp.Error)
		return 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	var response AgoricResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return 0, err
	}

	// loop through the validators voting powers
	for _, ele := range response.Result {
		val, _ := strconv.Atoi(ele.Tokens)
		votingPowers = append(votingPowers, int64(val))
	}

	// need to sort the powers in descending order since they are in random order
	sort.Slice(votingPowers, func(i, j int) bool { return votingPowers[i] > votingPowers[j] })

	totalVotingPower := utils.CalculateTotalVotingPower(votingPowers)
	fmt.Println("Total voting power:", totalVotingPower)

	// // now we're ready to calculate the nakomoto coefficient
	nakamotoCoefficient := utils.CalcNakamotoCoefficient(totalVotingPower, votingPowers)
	fmt.Println("The Nakamoto coefficient for agoric is", nakamotoCoefficient)

	return nakamotoCoefficient, nil
}
