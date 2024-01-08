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
	"strconv"
	"time"

	"github.com/xenowits/nakamoto-coefficient-calculator/core/utils"
)

const AptosValidatorsUrl = "https://fullnode.mainnet.aptoslabs.com/v1/accounts/0x1/resource/0x1::stake::ValidatorSet"

type AptosResponse struct {
	Data struct {
		ActiveValidators []struct {
			VotingPower string `json:"voting_power"`
		} `json:"active_validators"`
		TotalVotingPower string `json:"total_voting_power"`
	} `json:"data"`
}

func Aptos() (int, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFunc()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, AptosValidatorsUrl, nil)
	if err != nil {
		log.Println(err)
		return 0, errors.New("could not create get request for aptos")
	}

	resp, err := new(http.Client).Do(req)
	if err != nil {
		log.Println(err)
		return 0, errors.New("get request failed for aptos")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	var response AptosResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return 0, errors.New("could not unmarshal response for aptos")
	}

	expectedTotalVotingPower, err := strconv.ParseInt(response.Data.TotalVotingPower, 10, 64)
	if err != nil {
		return 0, errors.New("failed to convert total voting power to int64")
	}

	var votingPowers []big.Int

	for _, ele := range response.Data.ActiveValidators {
		val, _ := strconv.Atoi(ele.VotingPower)
		votingPowers = append(votingPowers, *big.NewInt(int64(val)))
	}

	calculatedTotalVotingPower := *utils.CalculateTotalVotingPowerBigNums(votingPowers)
	coefficient := utils.CalcNakamotoCoefficientBigNums(&calculatedTotalVotingPower, votingPowers)

	if expectedTotalVotingPower != calculatedTotalVotingPower.Int64() {
		fmt.Printf("Expected total voting power: %d\n", expectedTotalVotingPower)
		fmt.Printf("Calculated total voting power: %s\n", calculatedTotalVotingPower.String())
		return 0, fmt.Errorf("total voting power mismatch: expected %d != calculated %s", expectedTotalVotingPower, calculatedTotalVotingPower.String())
	}

	fmt.Printf("Total voting power: %s\n", calculatedTotalVotingPower.String())
	fmt.Printf("The Nakomoto coefficient for Aptos is %d\n", coefficient)
	return 1, nil
}
