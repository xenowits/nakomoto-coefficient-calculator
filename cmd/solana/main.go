package solana

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"os"
	"sort"
)

type Response []struct {
	Name         string `json:"name"`
	Account      string `json:"keybase_id"`
	Active_stake int    `json:"active_stake"`
	Delinquent   bool   `json:"delinquent"`
}

type ErrorResponse struct {
	Id      int    `json:"id"`
	Jsonrpc string `json:"jsonrpc"`
	Error   string `json:"error"`
}

func Solana() int {
	votingPowers := make([]big.Int, 0, 200)

	url := fmt.Sprintf("https://www.validators.app/api/v1/validators/mainnet.json")

	// Create a new request using http
	req, err := http.NewRequest("GET", url, nil)

	// add authorization header to the req
	// NOTE: get your own API_KEY from https://www.validators.app/api-documentation
	req.Header.Add("Token", os.Getenv("SOLANA_API_KEY"))

	// Send req using http Client
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		errBody, _ := ioutil.ReadAll(resp.Body)
		var errResp ErrorResponse
		json.Unmarshal(errBody, &errResp)
		log.Println(errResp.Error)
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Fatalln(err)
	}

	// loop through the validators voting powers
	for _, ele := range response {
		votingPowers = append(votingPowers, *big.NewInt(int64(ele.Active_stake)))
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
	fmt.Println("Total voting power:", new(big.Float).SetInt(totalVotingPower))

	// now we're ready to calculate the nakomoto coefficient
	nakamotoCoefficient := calcNakamotoCoefficient(totalVotingPower, votingPowers)
	fmt.Println("The Nakamoto coefficient for Solana is", nakamotoCoefficient)

	return nakamotoCoefficient
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
