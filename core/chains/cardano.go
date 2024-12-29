package chains

import (
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"sort"
	
	utils "github.com/xenowits/nakamoto-coefficient-calculator/core/utils"
)

type CardanoResponse struct {
	Label string  `json:"label"`
	Class string  `json:"class"`
	Epoch int     `json:"epoch"`
	Stake float64 `json:"stake"`
}

func Cardano() (int, error) {
    url := "https://www.balanceanalytics.io/api/mavdata.json"

    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        log.Println("Error creating request:", err)
        return 0, err
    }

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        log.Println("Error making request:", err)
        return 0, err
    }
    defer resp.Body.Close()

    var responseData struct {
        ApiData []CardanoResponse `json:"api_data"`
    }
    err = json.NewDecoder(resp.Body).Decode(&responseData)
    if err != nil {
        log.Println("Error decoding JSON:", err)
        return 0, err
    }

    var votingPowers []big.Int
    for _, data := range responseData.ApiData {
        stakeInt := big.NewInt(int64(data.Stake))
        votingPowers = append(votingPowers, *stakeInt)
    }

	// need to sort the powers in descending order since they are in random order
	sort.Slice(votingPowers, func(i, j int) bool {
        return votingPowers[i].Cmp(&votingPowers[j]) > 0
    })

	// Calculate total voting power
	totalVotingPower := utils.CalculateTotalVotingPowerBigNums(votingPowers)

	// Calculate Nakamoto coefficient
	nakamotoCoefficient := utils.CalcNakamotoCoefficientBigNums51(totalVotingPower, votingPowers)

	fmt.Println("The total voting power for Cardano is: ", totalVotingPower)
	fmt.Println("The Nakamoto coefficient for Cardano is: ", nakamotoCoefficient)

	// Return Nakamoto coefficient
	return nakamotoCoefficient, nil
}
