package chains

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strconv"

	utils "github.com/xenowits/nakamoto-coefficient-calculator/core/utils"
)

type RegenResponse struct {
	Height string `json:"height"`
	Result []struct {
		OperatorAddress string `json:"operator_address"`
		ConsensusPubkey struct {
			Type  string `json:"type"`
			Value string `json:"value"`
		} `json:"consensus_pubkey"`
		Tokens      string `json:"tokens"`
		Description struct {
			Moniker         string `json:"moniker"`
			Identity        string `json:"identity"`
			Website         string `json:"website"`
			SecurityContact string `json:"security_contact"`
			Details         string `json:"details"`
		} `json:"description"`
	} `json:"result"`
}

type RegenErrorResponse struct {
	Id      int    `json:"id"`
	Jsonrpc string `json:"jsonrpc"`
	Error   string `json:"error"`
}

func Regen() (int, error) {
	votingPowers := make([]int64, 0, 200)

	url := fmt.Sprintf("https://lcd-regen.keplr.app/staking/validators")
	resp, err := http.Get(url)
	if err != nil {
		errBody, _ := ioutil.ReadAll(resp.Body)
		var errResp RegenErrorResponse
		json.Unmarshal(errBody, &errResp)
		log.Println(errResp.Error)
		return -1, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return -1, err
	}

	var response RegenResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return -1, err
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
	fmt.Println("The Nakamoto coefficient for regen network is", nakamotoCoefficient)

	return nakamotoCoefficient, nil
}
