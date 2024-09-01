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
	validatorDataURL := "https://proxy.atomscan.com/cosmoshub-lcd/cosmos/staking/v1beta1/validators?page.offset=1&pagination.limit=500&status=BOND_STATUS_BONDED"
	stakingPoolURL := "https://proxy.atomscan.com/cosmoshub-lcd/cosmos/staking/v1beta1/pool"

	return FetchCosmosSDKNakaCoeff("cosmos", validatorDataURL, stakingPoolURL)
}

type cosmosValidatorData struct {
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

type cosmosStakingPoolData struct {
	Pool struct {
		NotBondedTokens string `json:"not_bonded_tokens"`
		BondedTokens    string `json:"bonded_tokens"`
	} `json:"pool"`
}

// fetchCosmosSDKNakaCoeff returns the nakamoto coefficient for a given cosmos SDK-based chain through REST API.
func FetchCosmosSDKNakaCoeff(chainName, validatorURL, poolURL string) (int, error) {
	var (
		votingPowers []big.Int
		validators   cosmosValidatorData
		pool         cosmosStakingPoolData
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

// Fetches data on active validator set
func fetchValidatorData(url string) (cosmosValidatorData, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFunc()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.Println(err)
		return cosmosValidatorData{}, errors.New("create get request for cosmos validators")
	}

	resp, err := new(http.Client).Do(req)
	if err != nil {
		log.Println(err)
		return cosmosValidatorData{}, errors.New("get request unsuccessful for cosmos validators")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return cosmosValidatorData{}, err
	}

	var response cosmosValidatorData
	err = json.Unmarshal(body, &response)
	if err != nil {
		return cosmosValidatorData{}, err
	}

	return response, nil
}

// Fetches staking pool data incl bonded and not_bonded tokens
func fetchStakingPoolData(url string) (cosmosStakingPoolData, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFunc()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.Println(err)
		return cosmosStakingPoolData{}, errors.New("create get request for cosmos pool")
	}

	resp, err := new(http.Client).Do(req)
	if err != nil {
		log.Println(err)
		return cosmosStakingPoolData{}, errors.New("get request unsuccessful for cosmos pool")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return cosmosStakingPoolData{}, err
	}

	var response cosmosStakingPoolData
	err = json.Unmarshal(body, &response)
	if err != nil {
		return cosmosStakingPoolData{}, err
	}

	return response, nil
}
