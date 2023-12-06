package chains

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	utils "github.com/xenowits/nakamoto-coefficient-calculator/core/utils"
)

type ApiResponse struct {
	LastUpdated  string    `json:"last_updated"`
	Validators   []int64   `json:"active_validator_balances"`
}

type ApiErrorResponse struct {
	Code         int       `json:"code"`
	Message      string    `json:"message"`
}

func Pulsechain() (int, error) {
	url := fmt.Sprintf("https://api.korkey.tech/pulsechain/validator_data.json")
	resp, err := http.Get(url)
	if err != nil {
		errBody, _ := ioutil.ReadAll(resp.Body)
		var errResp ApiErrorResponse
		
		errr := json.Unmarshal(errBody, &errResp)
		if errr != nil {
			return 0, errr
		}

		return 0, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	var response ApiResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return 0, err
	}

	// break if no entries
	if len(response.Validators) == 0 {
		return 0, err
	}

	totalVotingPower := utils.CalculateTotalVotingPower(response.Validators)
	fmt.Println("Total voting power:", totalVotingPower)

	// Now we're ready to calculate the nakamoto coefficient
	nakamotoCoefficient := utils.CalcNakamotoCoefficient(totalVotingPower, response.Validators)
	fmt.Println("The Nakamoto coefficient for Pulsechain is", nakamotoCoefficient)

	return nakamotoCoefficient, nil
}
