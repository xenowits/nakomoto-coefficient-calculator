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

func Cosmos() (int, error) {
	validatorURL := "https://proxy.atomscan.com/cosmoshub-lcd/cosmos/staking/v1beta1/validators?page.offset=1&pagination.limit=500&status=BOND_STATUS_BONDED"
	poolURL := "https://proxy.atomscan.com/cosmoshub-lcd/cosmos/staking/v1beta1/pool"

	return FetchCosmosSDKNakaCoeff("cosmos", validatorURL, poolURL)
}

type cosmosValidatorResponse struct {
	Validators []struct {
		OperatorAddress string `json:"operator_address"`
		ConsensusPubkey struct {
			Type string `json:"@type"`
			Key  string `json:"key"`
		} `json:"consensus_pubkey"`
		Jailed          bool   `json:"jailed"`
		Status          string `json:"status"`
		Tokens          string `json:"tokens"`
		DelegatorShares string `json:"delegator_shares"`
	} `json:"validators"`
}

type cosmosPoolResponse struct {
	Pool struct {
		NotBondedTokens string `json:"not_bonded_tokens"`
		BondedTokens    string `json:"bonded_tokens"`
	} `json:"pool"`
}

// fetchCosmosSDKNakaCoeff returns the nakamoto coefficient for a given cosmos SDK-based chain through REST API.
func FetchCosmosSDKNakaCoeff(chainName, validatorURL, poolURL string) (int, error) {
	var (
		votingPowers []big.Int
		validators   cosmosValidatorResponse
		pool         cosmosPoolResponse
		err          error
	)

	log.Printf("Fetching data for %s", chainName)

	// Fetch the validator data
	validators, err = fetchValidatorData(validatorURL)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch validator data for %s: %w", chainName, err)
	}

	// Fetch the staking pool data to get the total bonded tokens
	pool, err = fetchStakingPoolData(poolURL)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch pool data for %s: %w", chainName, err)
	}

	// Convert the bonded tokens from the pool response
	totalVotingPower, ok := new(big.Int).SetString(pool.Pool.BondedTokens, 10)
	if !ok {
		return 0, errors.New("failed to convert bonded tokens to big.Int")
	}

	// Loop through the validators' voting powers
	for _, ele := range validators.Validators {
		if ele.Status != BONDED {
			continue
		}

		val, err := strconv.Atoi(ele.Tokens)
		if err != nil {
			log.Printf("Error parsing token value for %s: %s", chainName, ele.Tokens)
			continue
		}
		votingPowers = append(votingPowers, *big.NewInt(int64(val)))
	}

	// Summarize voting powers for logging
	log.Printf("Voting powers for %s: %d validators with a total voting power of %s", chainName, len(votingPowers), totalVotingPower.String())

	if len(votingPowers) == 0 {
		return 0, fmt.Errorf("no valid voting powers found for %s", chainName)
	}

	// Sort the powers in descending order since they may be in random order
	sort.Slice(votingPowers, func(i, j int) bool {
		return votingPowers[i].Cmp(&votingPowers[j]) > 0
	})

	// Calculate the Nakamoto coefficient
	nakamotoCoefficient := utils.CalcNakamotoCoefficientBigNums(totalVotingPower, votingPowers)
	log.Printf("The Nakamoto coefficient for %s is %d", chainName, nakamotoCoefficient)

	return nakamotoCoefficient, nil
}

// Fetches
func fetchValidatorData(url string) (cosmosValidatorResponse, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFunc()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.Println(err)
		return cosmosValidatorResponse{}, errors.New("create get request for cosmos validators")
	}

	resp, err := new(http.Client).Do(req)
	if err != nil {
		log.Println(err)
		return cosmosValidatorResponse{}, errors.New("get request unsuccessful for cosmos validators")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return cosmosValidatorResponse{}, err
	}

	var response cosmosValidatorResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return cosmosValidatorResponse{}, err
	}

	return response, nil
}

func fetchStakingPoolData(url string) (cosmosPoolResponse, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFunc()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.Println(err)
		return cosmosPoolResponse{}, errors.New("create get request for cosmos pool")
	}

	resp, err := new(http.Client).Do(req)
	if err != nil {
		log.Println(err)
		return cosmosPoolResponse{}, errors.New("get request unsuccessful for cosmos pool")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return cosmosPoolResponse{}, err
	}

	var response cosmosPoolResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return cosmosPoolResponse{}, err
	}

	return response, nil
}
