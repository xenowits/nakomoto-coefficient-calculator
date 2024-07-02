package chains

import (
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"

	utils "github.com/xenowits/nakamoto-coefficient-calculator/core/utils"
)

type ChartData struct {
    AvgStake      float64 `json:"avgstake"`
    DelegateCount int     `json:"delegatecount"`
    Epoch         int     `json:"epoch"`
    Label         string  `json:"label"`
    Leverage      string  `json:"leverage"`
    MavGroup      string  `json:"mavgroup"`
    Pledge        float64 `json:"pledge"`
    PoolCount     int     `json:"poolcount"`
    PrctStake     float64 `json:"prctstake"`
    Stake         float64 `json:"stake"`
}

func Cardano() (int, error) {
    url := "https://api.balanceanalytics.io/rpc/pool_group_stake_donut"

    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        log.Println("Error creating request:", err)
        return 0, err
    }

    // This is the PUBLIC_BALANCE_API_TOKEN taken from https://www.balanceanalytics.io/chartboards/donut_shop
    req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJyb2xlIjoid2ViYXBwX3VzZXIifQ.eg3Zb9ZduibYJr1pgUrfqy4PFhkVU1uO_F9gFPBZnBI")
	req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        log.Println("Error making request:", err)
        return 0, err
    }
    defer resp.Body.Close()

    var responseData []struct {
        Chartdata []ChartData `json:"chartdata"`
    }
    err = json.NewDecoder(resp.Body).Decode(&responseData)
    if err != nil {
        log.Println("Error decoding JSON:", err)
        return 0, err
    }

	var votingPowers []big.Int
		for _, data := range responseData {
			for _, chartData := range data.Chartdata {
			stakeInt := big.NewInt(int64(chartData.Stake))
			votingPowers = append(votingPowers, *stakeInt)
		}
	}

	// Calculate total voting power
	totalVotingPower := utils.CalculateTotalVotingPowerBigNums(votingPowers)

	// Calculate Nakamoto coefficient
	nakamotoCoefficient := utils.CalcNakamotoCoefficientBigNums51(totalVotingPower, votingPowers)

	fmt.Println("Total voting power:", totalVotingPower)
	fmt.Println("The Nakamoto coefficient for Cardano is", nakamotoCoefficient)

	// Return Nakamoto coefficient
	return nakamotoCoefficient, nil
}
