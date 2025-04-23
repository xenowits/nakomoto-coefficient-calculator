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

type AvailResponse struct {
	Data struct {
		List []struct {
			BondedTotal string `json:"bonded_total"`
		} `json:"list"`
	} `json:"data"`
}

func Avail() (int, error) {
	var votingPowers []*big.Int
	ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFunc()

	url := "https://avail.api.subscan.io/api/scan/staking/validators"
	payload := bytes.NewBuffer([]byte(`{"order":"desc", "order_field":"bonded_total","row": 0,"page": 0}`))

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, payload)
	if err != nil {
		log.Println(err)
		return 0, errors.New("create post request for avail")
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

	var response AvailResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return 0, err
	}

	// Loop through the validators bonded amounts
	for _, ele := range response.Data.List {
		bondedTotal := new(big.Int)
		_, ok := bondedTotal.SetString(ele.BondedTotal, 10)
		if !ok {
			log.Println("Error parsing bonded total:", ele.BondedTotal)
			continue
		}
		
		votingPowers = append(votingPowers, bondedTotal)
	}

	// Sort the voting powers in descending order to ensure they're in correct order.
	sort.Slice(votingPowers, func(i, j int) bool { return votingPowers[i].Cmp(votingPowers[j]) > 0 })

	totalVotingPower := utils.CalculateTotalVotingPowerBigInt(votingPowers)
	fmt.Println("Total voting power:", totalVotingPower)

	// Calculate the nakamoto coefficient
	nakamotoCoefficient := utils.CalcNakamotoCoefficientBigInt(totalVotingPower, votingPowers)
	fmt.Println("The Nakamoto coefficient for Avail is", nakamotoCoefficient)

	return nakamotoCoefficient, nil
}
