package chains

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"
	"time"
)

type MinaResponse struct {
	Content []struct {
		Pk             string  `json:"pk"`
		Name           string  `json:"name"`
		StakePercent   float64 `json:"stakePercent"`
		CanonicalBlock int     `json:"canonicalBlock"`
		SocialTelegram string  `json:"socialTelegram"`
	}
	TotalPages    int `json:"totalPages"`
	TotalElements int `json:"totalElements"`
}

type MinaErrorResponse struct {
	Id      int    `json:"id"`
	Jsonrpc string `json:"jsonrpc"`
	Error   string `json:"error"`
}

func reverse(numbers []float64) {
	for i, j := 0, len(numbers)-1; i < j; i, j = i+1, j-1 {
		numbers[i], numbers[j] = numbers[j], numbers[i]
	}
}

func Mina() (int, error) {
	var votingPowers []float64
	var totalStake float64
	pageNo, entriesPerPage := 0, 50
	url := ""
	ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFunc()
	for true {
		// Check the most active url in the network logs here: https://mina.staketab.com/validators/stake
		// Sometimes it changes, like once it changed from mina.staketab.com to t-mina.staketab.com
		// Once, it was https://mina.staketab.com:8181/api/validator/all/
		url = fmt.Sprintf("https://minascan.io/mainnet/api/api/validators/?page=%d&size=%d&sortBy=amount_staked&type=active&findStr=&orderBy=DESC", pageNo, entriesPerPage)
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			log.Println(err)
			return 0, errors.New("create get request for mina")
		}

		resp, err := new(http.Client).Do(req)
		if err != nil {
			log.Println(err)
			return 0, errors.New("get request unsuccessful")
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return 0, err
		}

		var response MinaResponse
		err = json.Unmarshal(body, &response)
		if err != nil {
			return 0, err
		}

		// Break if no content or all pages have been fetched
		if len(response.Content) == 0 || pageNo >= response.TotalPages {
			break
		}

		// loop through the validators voting powers
		for _, ele := range response.Content {
			votingPowers = append(votingPowers, ele.StakePercent)
			// Accumulate total stake of all validators
			totalStake += ele.StakePercent
		}

		// increment counters
		pageNo += 1
	}

	// Sort voting powers in descending order
	sort.Slice(votingPowers, func(i, j int) bool {
		return votingPowers[i] > votingPowers[j]
	})

	// now we're ready to calculate the Nakamoto coefficient
	nakamotoCoefficient := calcNakamotoCoefficientForMina(votingPowers, totalStake)
	fmt.Println("The Nakamoto coefficient for Mina is", nakamotoCoefficient)

	return nakamotoCoefficient, nil
}

func calcNakamotoCoefficientForMina(votingPowers []float64, totalStake float64) int {
	// Calculate 50% threshold dynamically based on total stake
	var cumulativePercent float64 = 0.00
	thresholdPercent := totalStake * 0.50
	nakamotoCoefficient := 0

	for _, vpp := range votingPowers {
		// since this is the actual voting percentage, no need to multiply with 100
		cumulativePercent += vpp
		nakamotoCoefficient += 1
		// Check if cumulative voting power exceeds 50% of total stake
		if cumulativePercent >= thresholdPercent {
			break
		}
	}
	return nakamotoCoefficient
}
