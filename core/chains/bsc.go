package chains

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"

	utils "github.com/xenowits/nakamoto-coefficient-calculator/core/utils"
)

type Request struct {
	height   int
	page     int
	per_page int
}

type BscResponse struct {
	Code int `json:"code"`
	Data struct {
		Total      int `json:"total"`
		Validators []struct {
			OperatorAddress string `json:"operatorAddress"`
			Moniker         string `json:"moniker"`
			TotalStaked     string `json:"totalStaked"`
		} `json:"validators"`
	} `json:"data"`
}

type BscErrorResponse struct {
	Id      int    `json:"id"`
	Jsonrpc string `json:"jsonrpc"`
	Error   string `json:"error"`
}

// https://api.bnbchain.org/bnb-staking/v1/validator/all?limit=100&offset=0
func BSC() (int, error) {
	totalVotingPower := int64(0)
	votingPowers := make([]int64, 0, 200)
	pageLimit, pageOffset := 50, 0
	url := ""
	for true {
		url = fmt.Sprintf("https://api.bnbchain.org/bnb-staking/v1/validator/all?limit=%d&offset=%d", pageLimit, pageOffset)
		resp, err := http.Get(url)
		if err != nil {
			errBody, _ := ioutil.ReadAll(resp.Body)
			var errResp BscErrorResponse
			json.Unmarshal(errBody, &errResp)
			log.Println(errResp.Error)
			return 0, err
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return 0, err
		}

		var response BscResponse
		err = json.Unmarshal(body, &response)
		if err != nil {
			return 0, err
		}

		// break if no more entries left
		if len(response.Data.Validators) == 0 {
			break
		}

		// loop through the validators voting power proportions
		for _, ele := range response.Data.Validators {
			totalStaked := ele.TotalStaked
			wei, _ := new(big.Int).SetString(totalStaked, 10)
			votingPowers = append(votingPowers, weiToEther(wei))
			totalVotingPower += weiToEther(wei)
		}

		// increment counters
		pageOffset += pageLimit
	}

	// NOTE: NO need to calculate total voting power as the API response already
	// sends us proportional shares of each validator

	// now we're ready to calculate the nakomoto coefficient
	nakamotoCoefficient := utils.CalcNakamotoCoefficient(totalVotingPower, votingPowers)
	fmt.Println("The Nakamoto coefficient for BNB Smart Chain is", nakamotoCoefficient)

	return nakamotoCoefficient, nil
}

func weiToEther(wei *big.Int) int64 {
	f := new(big.Float)
	f.SetPrec(236) //  IEEE 754 octuple-precision binary floating-point format: binary256
	f.SetMode(big.ToNearestEven)
	fWei := new(big.Float)
	fWei.SetPrec(236) //  IEEE 754 octuple-precision binary floating-point format: binary256
	fWei.SetMode(big.ToNearestEven)
	result, _ := f.Quo(fWei.SetInt(wei), big.NewFloat(1e18)).Int64()
	return result
}
