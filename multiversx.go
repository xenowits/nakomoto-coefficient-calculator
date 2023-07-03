package chains

import (
	//"bytes"
	"encoding/json"
	"fmt"
	"github.com/xenowits/nakamoto-coefficient-calculator/core/utils"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"sort"
)

type ElrondResponse []struct {
		Stake                      string  `json:"stake"`
}

func Elrond() (int, error) {
	votingPowers := make([]big.Int, 0, 1024)

        url := fmt.Sprintf("https://api.elrond.com/nodes?type=validator")

	// Create a new POST request using http
	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	//req.Header.Set("Content-Type", "application/json")

	// Send req using http Client
	//client := &http.Client{}
	//resp, err := client.Do(req)

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
		return 0, err
	}

	body, err := ioutil.ReadAll(resp.Body)
        //fmt.Println(string(body))
        if err != nil {
		return 0, err
	}

	var response ElrondResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return 0, err
	}
	//fmt.Println((response))

	// loop through the validators voting powers
	for _, ele := range response {
		n, ok := new(big.Int).SetString(ele.Stake, 10)
		if !ok {
			return 0, fmt.Errorf("failed to parse string %s", ele.Stake)
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
	fmt.Println("The Nakamoto coefficient for Elrond protocol is", nakamotoCoefficient)

	return nakamotoCoefficient, nil
}
