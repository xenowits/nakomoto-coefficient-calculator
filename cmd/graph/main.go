package graph

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"sort"
)

type Response struct {
	Data struct {
		Indexers []struct {
			Id           string `json:"id"`
			StakedTokens string `json:"stakedTokens"`
		} `json:"indexers"`
	} `json:"data"`
}

type ErrorResponse struct {
	Id      int    `json:"id"`
	Jsonrpc string `json:"jsonrpc"`
	Error   string `json:"error"`
}

func Graph() (int, error) {
	votingPowers := make([]big.Int, 0, 1000)

	url := fmt.Sprintf("https://gateway.thegraph.com/network")
	jsonReqData := []byte(`{"query":"{ indexers (first: 1000) { id stakedTokens } }","variables":{}}`)

	// Create a new request using http
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonReqData))
	req.Header.Add("Content-Type", "application/json")

	// Send req using http Client
	client := &http.Client{}
	resp, err := client.Do(req)

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
	for _, ele := range response.Data.Indexers {
		n, ok := new(big.Int).SetString(ele.StakedTokens, 10)
		if !ok {
			log.Fatalln("Couldn't parse string", ele.StakedTokens)
		} else {
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

	totalVotingPower := calculateTotalVotingPower(votingPowers)
	fmt.Println("Total voting power:", totalVotingPower)

	// now we're ready to calculate the nakomoto coefficient
	nakamotoCoefficient := calcNakamotoCoefficient(totalVotingPower, votingPowers)
	fmt.Println("The Nakamoto coefficient for graph protocol is", nakamotoCoefficient)

	return nakamotoCoefficient, nil
}

func calculateTotalVotingPower(votingPowers []big.Int) *big.Int {
	total := big.NewInt(0)
	for _, vp := range votingPowers {
		total = new(big.Int).Add(total, &vp)
	}
	return total
}

func calcNakamotoCoefficient(totalVotingPower *big.Int, votingPowers []big.Int) int {
	thresholdPercent := big.NewFloat(0.33)
	thresholdVal := new(big.Float).Mul(new(big.Float).SetInt(totalVotingPower), thresholdPercent)
	cumulativeVal := big.NewFloat(0.00)
	nakamotoCoefficient := 0

	for _, vp := range votingPowers {
		z := new(big.Float).Add(cumulativeVal, new(big.Float).SetInt(&vp))
		cumulativeVal = z
		nakamotoCoefficient += 1
		if cumulativeVal.Cmp(thresholdVal) == +1 {
			break
		}
	}

	return nakamotoCoefficient
}