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

	utils "github.com/xenowits/nakamoto-coefficient-calculator/core/utils"
)

type AlgorandValidator struct {
	Address         string  `json:"address" `           // validator's address
	StakeMicroAlgo  uint64  `json:"stake_micro_algo" `  // stake in micro Algos
	StakeAlgo       float64 `json:"stake_algo" `        // stake in Algos
	StakeFraction   float64 `json:"stake_fraction" `    // fraction of total stake
	RewardsEligible bool    `json:"rewards_eligible" `  // is validator elifible for rewards
	ProposalsDaily  float64 `json:"proposals_daily" `   // expected numer of proposals daily
	Keytype         string  `json:"keytype" `           // signature type
	AsOfRound       uint64  `json:"as_of_round" `       // data valid as of round
	LastVotingRound uint64  `json:"last_voting_round" ` // participation key expires at this round
	ExpiresInDays   float64 `json:"expires_in_days" `   // part key expires in N days
}

type AlgorandResponse []AlgorandValidator

func Algorand() (int, error) {
	var votingPowers []int64
	ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFunc()

	// https://afmetrics.api.nodely.io/v1/api-docs/
	url := "https://afmetrics.api.nodely.io/v1/realtime/participation/validators"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.Println(err)
		return 0, errors.New("create get request for Algorand")
	}

	resp, err := new(http.Client).Do(req)
	if err != nil {
		log.Println(err)
		return 0, errors.New("get request unsuccessful")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}
	resp.Body.Close()

	var response AlgorandResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return 0, err
	}

	// Loop through the validators staked amounts
	for _, val := range response {
		votingPowers = append(votingPowers, int64(val.StakeMicroAlgo))
	}

	// Sort the voting powers in descending order since they maybe in random order.
	sort.Slice(votingPowers, func(i, j int) bool { return votingPowers[i] > votingPowers[j] })

	totalVotingPower := utils.CalculateTotalVotingPower(votingPowers)
	fmt.Println("Total voting power:", totalVotingPower)

	// Now we're ready to calculate the nakamoto coefficient
	nakamotoCoefficient := utils.CalcNakamotoCoefficient(totalVotingPower, votingPowers)
	fmt.Println("The Nakamoto coefficient for Algorand is", nakamotoCoefficient)

	return nakamotoCoefficient, nil
}
