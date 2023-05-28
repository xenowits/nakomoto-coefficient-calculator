package chains

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type Eth2Response struct {
	Total int `json:"total"`
	Data  []struct {
		NetworkPenetration float64 `json:"networkPenetration"`
	} `json:"data"`
}

// https://api.rated.network/docs#/
// the call used is for up to 100 entities but there are currently only 30
// should update if operator count grows beyond 100
//
// voting powers were not included as the api does not provide this metric
// and the ui does not display it
func Eth2() (int, error) {
	var (
		controllingPenetration = .33
		cumulativePenetration  = 0.0
		nakamotoCoefficient    = 1
	)

	ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFunc()

	url := fmt.Sprintf("https://api.rated.network/v0/eth/operators?window=all&idType=entity&size=100")
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.Println(err)
		return 0, errors.New("create get request for ethereum")
	}

	// Add authorization header to the request
	// NOTE: You can get your own API_KEY from https://bit.ly/ratedAPIkeys
	authToken := fmt.Sprintf("Bearer %s", os.Getenv("RATED_API_KEY"))
	req.Header.Add("Authorization", authToken)

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

	var response Eth2Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		return 0, err
	}

	// loop through operators and pull their network penetration
	// increments NC when the controllingPenetration is surpassed
	for _, ele := range response.Data {
		cumulativePenetration += ele.NetworkPenetration
		if cumulativePenetration < controllingPenetration {
			nakamotoCoefficient += 1
		}
	}

	fmt.Println("The Nakamoto coefficient for Eth2 is", nakamotoCoefficient)

	return nakamotoCoefficient, nil
}
