package chains

import (
	"encoding/json"
	"fmt"
	"github.com/xenowits/nakamoto-coefficient-calculator/core/utils"
	"io"
	"log"
	"math/big"
	"net/http"
	"sort"
)

type RegenResponse struct {
	Data []struct {
		Tokens int64 `json:"tokens"`
	} `json:"data"`
}

func Regen() (int, error) {
	const url = "https://api.regen.aneka.io/validators/details/all"

	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	var response RegenResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return 0, err
	}

	// Loop through the validators voting powers.
	var votingPowers []big.Int
	for _, ele := range response.Data {
		votingPowers = append(votingPowers, *big.NewInt(ele.Tokens))
	}

	// Need to sort the powers in descending order since they maybe in random order.
	sort.Slice(votingPowers, func(i, j int) bool {
		res := (&votingPowers[i]).Cmp(&votingPowers[j])
		if res == 1 {
			return true
		}
		return false
	})

	totalVotingPower := utils.CalculateTotalVotingPowerBigNums(votingPowers)
	fmt.Println("Total voting power:", new(big.Float).SetInt(totalVotingPower))

	// now we're ready to calculate the nakamoto coefficient
	nakamotoCoefficient := utils.CalcNakamotoCoefficientBigNums(totalVotingPower, votingPowers)
	fmt.Println("The Nakamoto coefficient for regen network is", nakamotoCoefficient)

	return nakamotoCoefficient, nil
}
