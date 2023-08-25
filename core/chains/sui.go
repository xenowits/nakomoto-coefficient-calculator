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
	"strconv"
	"time"

	utils "github.com/xenowits/nakamoto-coefficient-calculator/core/utils"
)

type SuiResponse struct {
	Result struct {
		ActiveValidators []struct {
			VotingPower string `json:"votingPower"`
		} `json:"activeValidators"`
	} `json:"result"`
}

type rawBody struct {
	JSONRPC string        `json:"jsonrpc"`
	ID      int           `json:"id"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
}

func Sui() (int, error) {
	request := rawBody{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "suix_getLatestSuiSystemState",
		Params:  []interface{}{},
	}

	baseURL := "https://fullnode.mainnet.sui.io"

	return fetchDataSUI("sui", baseURL, request)
}

func fetchDataSUI(chainName string, url string, request rawBody) (int, error) {
	var votingPowers []big.Int

	response, err := fetchData(url, request)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch data for %s: %w", chainName, err)
	}

	// Loop through the validators voting powers.
	// Sui has voting power indicator for each validator.
	// Total is 10000==100%
	for _, ele := range response.Result.ActiveValidators {
		votingPower, err := strconv.ParseInt(ele.VotingPower, 10, 64)
		if err != nil {
			log.Println(err)
		}

		votingPowers = append(votingPowers, *big.NewInt(votingPower))
	}

	sort.Slice(votingPowers, func(i, j int) bool {
		return votingPowers[i].Cmp(&votingPowers[j]) > 0
	})

	totalVotingPower := utils.CalculateTotalVotingPowerBigNums(votingPowers)
	fmt.Printf("Total voting power for %s: %s\n", chainName, new(big.Float).SetInt(totalVotingPower).String())

	// Now we're ready to calculate the nakomoto coefficient.
	nakamotoCoefficient := utils.CalcNakamotoCoefficientBigNums(totalVotingPower, votingPowers)
	fmt.Printf("The Nakamoto coefficient for %s is %d\n", chainName, nakamotoCoefficient)

	return nakamotoCoefficient, nil
}

func fetchData(url string, request rawBody) (SuiResponse, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFunc()

	requestBody, err := json.Marshal(request)
	if err != nil {
		log.Println(err)
		return SuiResponse{}, errors.New("failed to marshal request for sui")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(requestBody))
	if err != nil {
		log.Println(err)
		return SuiResponse{}, errors.New("create POST request for sui")
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return SuiResponse{}, errors.New("POST request unsuccessful for sui")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return SuiResponse{}, err
	}

	var response SuiResponse

	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Println(err)
		return SuiResponse{}, errors.New("failed to unmarshal response for sui")
	}

	return response, nil
}
