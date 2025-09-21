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
	"time"

	utils "github.com/xenowits/nakamoto-coefficient-calculator/core/utils"
)

type NamadaValidator struct {
	VotingPower string `json:"voting_power"`
}

type NamadaValidatorsResponse struct {
	Result struct {
		Validators []NamadaValidator `json:"validators"`
	} `json:"result"`
}

type NamadaTotalVotingPowerResponse struct {
	TotalVotingPower string `json:"totalVotingPower"`
}

func Namada() (int, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFunc()

	// Fetch validators
	validatorsURL := "https://namada-archive.tm.p2p.org/validators"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, validatorsURL, nil)
	if err != nil {
		log.Println(err)
		return 0, errors.New("create get request for namada validators")
	}
	resp, err := new(http.Client).Do(req)
	if err != nil {
		log.Println(err)
		return 0, errors.New("get request unsuccessful for namada validators")
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}
	var valResp NamadaValidatorsResponse
	err = json.Unmarshal(body, &valResp)
	if err != nil {
		return 0, err
	}

	// Fetch total voting power
	totalPowerURL := "https://api-namada-mainnet-indexer.tm.p2p.org/api/v1/pos/voting-power"
	req2, err := http.NewRequestWithContext(ctx, http.MethodGet, totalPowerURL, nil)
	if err != nil {
		log.Println(err)
		return 0, errors.New("create get request for namada total voting power")
	}
	resp2, err := new(http.Client).Do(req2)
	if err != nil {
		log.Println(err)
		return 0, errors.New("get request unsuccessful for namada total voting power")
	}
	defer resp2.Body.Close()
	body2, err := io.ReadAll(resp2.Body)
	if err != nil {
		return 0, err
	}
	var totalResp NamadaTotalVotingPowerResponse
	err = json.Unmarshal(body2, &totalResp)
	if err != nil {
		return 0, err
	}

	// Parse voting powers
	var votingPowers []*big.Int
	for _, v := range valResp.Result.Validators {
		vp := new(big.Int)
		_, ok := vp.SetString(v.VotingPower, 10)
		if !ok {
			log.Println("Error parsing validator voting power:", v.VotingPower)
			continue
		}
		votingPowers = append(votingPowers, vp)
	}

	// Sort voting powers descending
	sort.Slice(votingPowers, func(i, j int) bool { return votingPowers[i].Cmp(votingPowers[j]) > 0 })

	totalVotingPower := new(big.Int)
	_, ok := totalVotingPower.SetString(totalResp.TotalVotingPower, 10)
	if !ok {
		return 0, fmt.Errorf("error parsing total voting power: %s", totalResp.TotalVotingPower)
	}
	// Multiply by 10^6 as required by Namada
	multiplier := big.NewInt(1_000_000)
	totalVotingPower.Mul(totalVotingPower, multiplier)

	nakamotoCoefficient := utils.CalcNakamotoCoefficientBigInt(totalVotingPower, votingPowers)
	fmt.Println("The Nakamoto coefficient for Namada is", nakamotoCoefficient)

	return nakamotoCoefficient, nil
}
