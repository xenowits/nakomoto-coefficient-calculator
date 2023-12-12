package chains

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"
	"strconv"
	"time"

	utils "github.com/xenowits/nakamoto-coefficient-calculator/core/utils"
)

type PolkadotResponse struct {
	Data struct {
		List []struct {
			BondedTotal string `json:"bonded_total"`
		} `json:"list"`
	} `json:"data"`
}

func Polkadot() (int, error) {
	var votingPowers []int64
	ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFunc()

	url := "https://polkadot.api.subscan.io/api/scan/staking/validators"
	payload := bytes.NewBuffer([]byte(`{"order":"desc", "order_field":"bonded_total","row": 0,"page": 0}`))

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, payload)
	if err != nil {
		log.Println(err)
		return 0, errors.New("create post request for polkadot")
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

	var response PolkadotResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return 0, err
	}

	// Loop through the validators bonded amounts
	for _, ele := range response.Data.List {
		bondedTotal, err := strconv.ParseInt(ele.BondedTotal, 10, 64)
		if err != nil {
			log.Println(err)
			continue
		}
		
		votingPowers = append(votingPowers, bondedTotal)
	}

	// Sort the voting powers in descending order since they maybe in random order.
	sort.Slice(votingPowers, func(i, j int) bool { return votingPowers[i] > votingPowers[j] })

	totalVotingPower := utils.CalculateTotalVotingPower(votingPowers)
	fmt.Println("Total voting power:", totalVotingPower)

	// Now we're ready to calculate the nakamoto coefficient
	nakamotoCoefficient := utils.CalcNakamotoCoefficient(totalVotingPower, votingPowers)
	fmt.Println("The Nakamoto coefficient for Polkadot is", nakamotoCoefficient)

	return nakamotoCoefficient, nil
}
