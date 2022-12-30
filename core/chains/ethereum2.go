package chains

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Eth2Response struct {
	Page struct {
		From int `json:"from"`
		Size int `json:"size"`
	} `json:"page"`
	Total int    `json:"total"`
	Next  string `json:"next"`
	Data  []struct {
		Id                        string  `json:"id"`
		IdType                    string  `json:"idType"`
		TimeWindow                string  `json:"timeWindow"`
		ValidatorCount            int     `json:"validatorCount"`
		AvgCorrectness            float64 `json:"avgCorrectness"`
		AvgInclusionDelay         float64 `json:"avgInclusionDelay"`
		AvgUptime                 float64 `json:"avgUptime"`
		AvgValidatorEffectiveness float64 `json:"avgValidatorEffectiveness"`
		ClientPercentages         []struct {
			Client     string  `json:"client"`
			Percentage float64 `json:"percentage"`
		} `json:"clientPercentages"`
		NetworkPenetration float64 `json:"networkPenetration"`
	} `json:"data"`
}

type Eth2ErrorResponse struct {
	Detail string `json:"detail"`
}

// https://api.rated.network/docs#/
// the call used is for up to 100 entities but there are currently only 30
// should update if operator count grows beyond 100
//
// voting powers were not included as the api does not provide this metric
// and the ui does not display it
func Eth2() (int, error) {
	controllingPenetration := .33
	cumulativePenetration := 0.0
	nakamotoCoefficient := 1
	url := ""

	url = fmt.Sprintf("https://api.rated.network/v0/eth/operators?window=all&idType=entity&size=100")
	resp, err := http.Get(url)
	if err != nil {
		errBody, _ := ioutil.ReadAll(resp.Body)
		var errResp Eth2ErrorResponse
		json.Unmarshal(errBody, &errResp)
		return -1, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return -1, err
	}

	var response Eth2Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		return -1, err
	}

	// loop through operators and pull their network penetration
	// increments NC when the controllingPenetration is surpassed
	for _, ele := range response.Data {
		val := ele.NetworkPenetration
		cumulativePenetration += val
		if cumulativePenetration < controllingPenetration {
			nakamotoCoefficient += 1
		}
	}

	fmt.Println("The Nakamoto coeffiecient for Eth2 is", nakamotoCoefficient)

	return nakamotoCoefficient, nil
}
