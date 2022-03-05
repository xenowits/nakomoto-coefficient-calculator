package chains

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"math/big"
	"net/http"
	"sort"

	utils "github.com/xenowits/nakamoto-coefficient-calculator/core/utils"
)

type JunoResponse struct {
	Data struct {
		Validators []struct {
			ValidatorVotingPowers []struct {
				VotingPower uint64 `json:"votingPower"`
			} `json:"validatorVotingPowers"`
		} `json:"validator"`
	} `json:"data"`
}

type JunoErrorResponse struct {
	Id      int    `json:"id"`
	Jsonrpc string `json:"jsonrpc"`
	Error   string `json:"error"`
}

func Juno() (int, error) {
	votingPowers := make([]big.Int, 0, 1000)

	url := fmt.Sprintf("https://hasura.junoscan.com/v1/graphql")
	jsonReqData := "{\"query\":\"query { validator { validatorVotingPowers: validator_voting_powers(order_by: { voting_power: desc }) { votingPower: voting_power } }}\",\"variables\":{}}"

	// Create a new request using http
	req, err := http.NewRequest("POST", url, strings.NewReader(jsonReqData))
	if err != nil {
		return -1, err
	}
	req.Header.Add("Content-Type", "application/json")

	// Send req using http Client
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		errBody, _ := ioutil.ReadAll(resp.Body)
		var errResp JunoErrorResponse
		json.Unmarshal(errBody, &errResp)
		log.Println(errResp.Error)
		return -1, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return -1, err
	}

	var response JunoResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return -1, err
	}

	// loop through the validators voting powers
	for _, ele := range response.Data.Validators {
		if len(ele.ValidatorVotingPowers) > 0 {
			n := new(big.Int).SetUint64(ele.ValidatorVotingPowers[0].VotingPower)
			votingPowers = append(votingPowers, *n)
		}
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

	// now we're ready to calculate the nakomoto coefficient
	nakamotoCoefficient := utils.CalcNakamotoCoefficientBigNums(totalVotingPower, votingPowers)
	fmt.Println("The Nakamoto coefficient for juno is", nakamotoCoefficient)

	return nakamotoCoefficient, nil
}
