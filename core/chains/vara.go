package chains

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math/big"
	"net/http"
	"sort"
	"time"

	utils "github.com/xenowits/nakamoto-coefficient-calculator/core/utils"
)

type VaraResponse struct {
	Data struct {
		List []struct {
			BondedTotal string `json:"bonded_total"`
		} `json:"list"`
	} `json:"data"`
}

func Vara() (int, error) {
	var votingPowers []*big.Int
	ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFunc()

	url := "https://vara.api.subscan.io/api/scan/staking/validators"
	payload := bytes.NewBuffer([]byte(`{"order":"desc", "order_field":"bonded_total","row": 0,"page": 0}`))

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, payload)
	if err != nil {
		log.Println(err)
		return 0, errors.New("create post request for vara")
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := new(http.Client).Do(req)
	if err != nil {
		log.Println(err)
		return 0, errors.New("post request unsuccessful")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}
	resp.Body.Close()

	var response VaraResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Printf("Error unmarshalling Vara response: %v", err)
		return 0, err
	}

	for _, ele := range response.Data.List {
		bondedTotal := new(big.Int)
		_, success := bondedTotal.SetString(ele.BondedTotal, 10)
		if !success {
			log.Printf("Error parsing BondedTotal '%s' into big.Int", ele.BondedTotal)
			continue
		}

		votingPowers = append(votingPowers, bondedTotal)
	}

	// Sort the voting powers (slice of *big.Int) in descending order.
	sort.Slice(votingPowers, func(i, j int) bool {
		// Compare returns -1 if votingPowers[i] < votingPowers[j],
		// 0 if votingPowers[i] == votingPowers[j],
		// +1 if votingPowers[i] > votingPowers[j].
		// We want descending order, so return true if votingPowers[i] > votingPowers[j].
		return votingPowers[i].Cmp(votingPowers[j]) == 1
	})

	// Convert []*big.Int to []big.Int for utility functions
	votingPowersValues := make([]big.Int, len(votingPowers))
	for i, ptr := range votingPowers {
		if ptr != nil {
			votingPowersValues[i] = *ptr
		} else {
			votingPowersValues[i] = *big.NewInt(0)
			log.Printf("Warning: Found nil big.Int pointer at index %d in Vara votingPowers", i)
		}
	}

	totalVotingPower := utils.CalculateTotalVotingPowerBigNums(votingPowersValues)
	fmt.Println("Total voting power (Vara):", totalVotingPower.String())

	nakamotoCoefficient := utils.CalcNakamotoCoefficientBigNums(totalVotingPower, votingPowersValues)
	fmt.Println("The Nakamoto coefficient for Vara Network is", nakamotoCoefficient)

	return nakamotoCoefficient, nil
}
