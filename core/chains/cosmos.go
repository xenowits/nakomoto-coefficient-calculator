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

type cosmosResponse struct {
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
		Description     struct {
			Moniker         string `json:"moniker"`
			Identity        string `json:"identity"`
			Website         string `json:"website"`
			SecurityContact string `json:"security_contact"`
			Details         string `json:"details"`
		} `json:"description"`
		UnbondingHeight string `json:"unbonding_height"`
		UnbondingTime   string `json:"unbonding_time"`
		Commission      struct {
			CommissionRates struct {
				Rate          string `json:"rate"`
				MaxRate       string `json:"max_rate"`
				MaxChangeRate string `json:"max_change_rate"`
			} `json:"commission_rates"`
			UpdateTime string `json:"update_time"`
		} `json:"commission"`
		MinSelfDelegation       string   `json:"min_self_delegation"`
		UnbondingOnHoldRefCount string   `json:"unbonding_on_hold_ref_count"`
		UnbondingIds            []string `json:"unbonding_ids"`
		ValidatorBondShares     string   `json:"validator_bond_shares"`
		LiquidShares            string   `json:"liquid_shares"`
	} `json:"validators"`
}

type cosmosPoolResponse struct {
	Pool struct {
		NotBondedTokens string `json:"not_bonded_tokens"`
		BondedTokens    string `json:"bonded_tokens"`
	} `json:"pool"`
}

func FetchCosmosSDKNakaCoeff(chainName, validatorURL, poolURL string) (int, error) {
	var (
		votingPowers []big.Int
		validators   cosmosResponse
		pool         cosmosPoolResponse
		err          error
	)

	log.Printf("Fetching data for %s", chainName)

	// Fetch the validator data
	validators, err = fetchValidators(validatorURL)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch validator data for %s: %w", chainName, err)
	}

	// Fetch the pool data to get the total bonded tokens
	pool, err = fetchPool(poolURL)
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

	// Summarize the voting powers for logging
	log.Printf("Voting powers for %s: %d validators with a total voting power of %s", chainName, len(votingPowers), totalVotingPower.String())

	if len(votingPowers) == 0 {
		return 0, fmt.Errorf("no valid voting powers found for %s", chainName)
	}

	// Sort the powers in descending order since they may be in random order
	sort.Slice(votingPowers, func(i, j int) bool {
		return votingPowers[i].Cmp(&votingPowers[j]) > 0
	})

	// Now we're ready to calculate the nakamoto coefficient
	nakamotoCoefficient := utils.CalcNakamotoCoefficientBigNums(totalVotingPower, votingPowers)
	log.Printf("The Nakamoto coefficient for %s is %d", chainName, nakamotoCoefficient)

	return nakamotoCoefficient, nil
}

func fetchValidators(url string) (cosmosResponse, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFunc()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.Println(err)
		return cosmosResponse{}, errors.New("create get request for cosmos validators")
	}

	resp, err := new(http.Client).Do(req)
	if err != nil {
		log.Println(err)
		return cosmosResponse{}, errors.New("get request unsuccessful for cosmos validators")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return cosmosResponse{}, err
	}

	var response cosmosResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return cosmosResponse{}, err
	}

	return response, nil
}

func fetchPool(url string) (cosmosPoolResponse, error) {
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
