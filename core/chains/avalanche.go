package chains

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"sort"
	"strconv"

	utils "github.com/xenowits/nakamoto-coefficient-calculator/core/utils"
)

type AvalancheResponse struct {
	Jsonrpc string `json:"jsonrpc"`
	Id      int    `json:"id"`
	Result  struct {
		Validators []struct {
			StakeAmount string `json:"stakeAmount"`
			NodeId      string `json:"nodeID"`
		} `json:"validators"`
	} `json:"result"`
}

type AvalancheErrorResponse struct {
	Id      int    `json:"id"`
	Jsonrpc string `json:"jsonrpc"`
	Error   string `json:"error"`
}

// Avalanche calculates nakamoto coefficient for the `Avalanche` C-Chain
// In AVAX, stake amounts are already multiplied by 10^9
// So, we need to deal with big numbers here.
// Else, if we divide each value with 10^9, we have to deal with fractional numbers which is worse.
func Avalanche() (int, error) {
	votingPowers := make([]big.Int, 0)

	url := fmt.Sprintf("https://api.avax.network/ext/P")
	jsonReqData := []byte(`{"jsonrpc": "2.0","method": "platform.getCurrentValidators","params":{},"id":1}`)

	// Create a new POST request using http
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonReqData))
	req.Header.Set("Content-Type", "application/json")

	// Send request using http Client
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		errBody, _ := ioutil.ReadAll(resp.Body)
		var errResp AvalancheErrorResponse
		json.Unmarshal(errBody, &errResp)
		log.Println(errResp.Error)
		return -1, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return -1, err
	}

	var response AvalancheResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return -1, err
	}

	// loop through the validators voting powers
	for _, ele := range response.Result.Validators {
		val, _ := strconv.Atoi(ele.StakeAmount)
		votingPowers = append(votingPowers, *big.NewInt(int64(val)))
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

	// // now we're ready to calculate the nakomoto coefficient
	nakamotoCoefficient := utils.CalcNakamotoCoefficientBigNums(totalVotingPower, votingPowers)
	fmt.Println("The Nakamoto coefficient for Avalanche is", nakamotoCoefficient)

	return nakamotoCoefficient, nil
}
