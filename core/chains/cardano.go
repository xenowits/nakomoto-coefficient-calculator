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
	Code int      `json:"code"`
	Time string   `json:"time"`
	Msg  string   `json:"msg"`
	Data []CardanoPool `json:"data"`
}

type CardanoPool struct {
	PoolID         string  `json:"pool_id"`
	Name           string  `json:"name"`
	Stake          string  `json:"stake"`
	BlocksLifetime string  `json:"blocks_lifetime"`
	ROALifetime    string  `json:"roa_lifetime"`
	Pledge         string  `json:"pledge"`
	Delegators     string  `json:"delegators"`
	Saturation     float64 `json:"saturation"`
}

func Cardano() (int, error) {
	url := "https://js.cexplorer.io/api-static/pool/list.json"

	var stakeAmounts []big.Int

	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	defer resp.Body.Close()

	var response CardanoResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return 0, err
	}

	// Loop through the pools and extract stake amounts
	for _, pool := range response.Data {
		stakeInt, ok := new(big.Int).SetString(pool.Stake, 10)
		if !ok {
			log.Println("Error converting stake amount to big.Int")
			continue
		}
		stakeAmounts = append(stakeAmounts, *stakeInt)
	}

	// need to sort the stake amounts in descending order
	sort.Slice(stakeAmounts, func(i, j int) bool {
		return stakeAmounts[i].Cmp(&stakeAmounts[j]) == 1
	})

	totalStake := utils.CalculateTotalVotingPowerBigNums(stakeAmounts)
	fmt.Println("Total voting power:", new(big.Float).SetInt(totalStake))

	// now we're ready to calculate the nakamoto coefficient
	nakamotoCoefficient := utils.CalcNakamotoCoefficientBigNums(totalStake, stakeAmounts)
	fmt.Println("The Nakamoto coefficient for Cardano is", nakamotoCoefficient)

	return nakamotoCoefficient, nil
}
