package chains

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/xenowits/nakamoto-coefficient-calculator/core/utils"
)

type TerraResponse struct {
	Validators []struct {
		Voting_power      string `json:"voting_power"`
		Proposer_priority string `json:"proposer_priority"`
	} `json:"validators"`
}

// Terra returns the NC value for terra blockchain.
//
// https://fcd.terra.dev/swagger
func Terra() (int, error) {
	votingPowers := make([]int64, 0, 200)
	ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFunc()

	url := fmt.Sprintf("https://phoenix-lcd.terra.dev/validatorsets/latest")
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.Println(err)
		return 0, errors.New("create get request for terra")
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

	var response TerraResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return 0, err
	}

	// loop through the validators voting powers
	for _, ele := range response.Validators {
		val, _ := strconv.Atoi(ele.Voting_power)
		votingPowers = append(votingPowers, int64(val))
	}

	// No need to sort as the result is already in sorted in descending order
	totalVotingPower := utils.CalculateTotalVotingPower(votingPowers)
	fmt.Println("Total voting power:", totalVotingPower)

	// now we're ready to calculate the nakomoto coefficient
	nakamotoCoefficient := utils.CalcNakamotoCoefficient(totalVotingPower, votingPowers)
	fmt.Println("The Nakamoto coefficient for terra is", nakamotoCoefficient)

	return nakamotoCoefficient, nil
}
