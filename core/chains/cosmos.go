package chains

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math/big"
	"net/http"
	"sort"
	"strconv"
	"time"

	utils "github.com/xenowits/nakamoto-coefficient-calculator/core/utils"
)

const BONDED = "BOND_STATUS_BONDED"

type cosmosResponse struct {
	Validators []struct {
		Status string `json:"status"`
		Tokens string `json:"tokens"`
	} `json:"validators"`
}

func Cosmos() (int, error) {
	url := "https://proxy.atomscan.com/cosmoshub-lcd/cosmos/staking/v1beta1/validators?page.offset=1&pagination.limit=500&status=BOND_STATUS_BONDED"
	return fetchCosmosSDKNakamotoCoefficient("cosmos", url)
}

func fetchCosmosSDKNakamotoCoefficient(chainName, url string) (int, error) {
	var (
		votingPowers []big.Int
		response     cosmosResponse
		err          error
	)

	response, err = fetch(url)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch data for %s: %w", chainName, err)
	}

	// loop through the validators voting powers
	for _, ele := range response.Validators {
		if ele.Status != BONDED {
			continue
		}

		val, _ := strconv.Atoi(ele.Tokens)
		votingPowers = append(votingPowers, *big.NewInt(int64(val)))
	}

	// Sort the powers in descending order since they maybe in random order
	sort.Slice(votingPowers, func(i, j int) bool {
		res := (&votingPowers[i]).Cmp(&votingPowers[j])
		return res == 1
	})

	totalVotingPower := utils.CalculateTotalVotingPowerBigNums(votingPowers)
	fmt.Println("Total voting power:", totalVotingPower)

	// Now we're ready to calculate the nakamoto coefficient
	nakamotoCoefficient := utils.CalcNakamotoCoefficientBigNums(totalVotingPower, votingPowers)
	fmt.Printf("The Nakamoto coefficient for %s is %d\n", chainName, nakamotoCoefficient)

	return nakamotoCoefficient, nil
}
func fetch(url string) (cosmosResponse, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFunc()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.Println(err)
		return cosmosResponse{}, errors.New("create get request for cosmos")
	}

	resp, err := new(http.Client).Do(req)
	if err != nil {
		log.Println(err)
		return cosmosResponse{}, errors.New("get request unsuccessful for cosmos")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return cosmosResponse{}, err
	}

	var response cosmosResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return cosmosResponse{}, nil
	}

	return response, nil
}
