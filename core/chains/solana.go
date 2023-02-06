package chains

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/big"
	"net/http"
	"os"
	"sort"

	utils "github.com/xenowits/nakamoto-coefficient-calculator/core/utils"
)

type SolanaResponse []struct {
	Name         string `json:"name"`
	Account      string `json:"keybase_id"`
	Active_stake int64  `json:"active_stake"`
	Delinquent   bool   `json:"delinquent"`
}

func Solana() (int, error) {
	url := fmt.Sprintf("https://www.validators.app/api/v1/validators/mainnet.json")

	var votingPowers []big.Int

	// Create a new request using http
	req, err := http.NewRequest("GET", url, nil)

	// Add authorization header to the request
	// NOTE: You can get your own API_KEY from https://www.validators.app/api-documentation
	req.Header.Add("Token", os.Getenv("SOLANA_API_KEY"))

	// Send req using http Client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	var response SolanaResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return 0, err
	}

	// loop through the validators voting powers
	for _, ele := range response {
		votingPowers = append(votingPowers, *big.NewInt(ele.Active_stake))
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
	fmt.Println("Total voting power:", new(big.Float).SetInt(totalVotingPower))

	// now we're ready to calculate the nakomoto coefficient
	nakamotoCoefficient := utils.CalcNakamotoCoefficientBigNums(totalVotingPower, votingPowers)
	fmt.Println("The Nakamoto coefficient for Solana is", nakamotoCoefficient)

	return nakamotoCoefficient, nil
}
