package chains

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"

	utils "github.com/xenowits/nakamoto-coefficient-calculator/core/utils"
)

const BOND_STATUS_BONDED = "BOND_STATUS_BONDED"
const JunoValidatorsUrl = "https://validators.cosmos.directory/chains/juno"

type JunoValidators struct {
	Name       string `json:"name"`
	Validators []struct {
		Status string `json:""`
		Tokens int64  `json:",string"`
		Jailed bool   `json:""`
	} `json:"validators"`
}

func Juno() (int, error) {
	var votingPowers []int64

	resp, err := http.Get(JunoValidatorsUrl)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	var response JunoValidators
	err = json.Unmarshal(body, &response)
	if err != nil {
		return 0, err
	}

	// loop through the validators voting powers
	for _, ele := range response.Validators {
		if ele.Jailed || ele.Status != BOND_STATUS_BONDED {
			continue
		}

		votingPowers = append(votingPowers, ele.Tokens)
	}

	// need to sort the powers in descending order since they are in random order
	sort.Slice(votingPowers, func(i, j int) bool { return votingPowers[i] > votingPowers[j] })

	totalVotingPower := utils.CalculateTotalVotingPower(votingPowers)
	fmt.Println("Total voting power:", totalVotingPower)

	// now we're ready to calculate the nakomoto coefficient
	nakamotoCoefficient := utils.CalcNakamotoCoefficient(totalVotingPower, votingPowers)
	fmt.Println("The Nakamoto coefficient for juno is", nakamotoCoefficient)

	return nakamotoCoefficient, nil
}
