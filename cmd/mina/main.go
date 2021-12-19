package mina

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"

	utils "github.com/xenowits/nakamoto-coefficient-calculator/utils"
)

type Request struct {
	height   int
	page     int
	per_page int
}

type Response struct {
	Content []struct {
		Pk             string  `json:"pk"`
		Name           string  `json:"name"`
		StakePercent   float64 `json:"stakePercent"`
		CanonicalBlock int     `json:"canonicalBlock"`
		SocialTelegram string  `json:"socialTelegram"`
	}
	TotalPages    int `json:"totalPages"`
	TotalElements int `json:"totalElements"`
}

type ErrorResponse struct {
	Id      int    `json:"id"`
	Jsonrpc string `json:"jsonrpc"`
	Error   string `json:"error"`
}

func reverse(numbers []float64) {
	for i, j := 0, len(numbers)-1; i < j; i, j = i+1, j-1 {
		numbers[i], numbers[j] = numbers[j], numbers[i]
	}
}

func Mina() int {
	votingPowers := make([]float64, 0, 200)
	pageNo, entriesPerPage := 0, 50
	url := ""
	for true {
		url = fmt.Sprintf("https://mina.staketab.com:8181/api/validator/all/?page=%d&size=%d&sortBy=canonical_block&findStr=&orderBy=DESC", pageNo, entriesPerPage)
		resp, err := http.Get(url)
		if err != nil {
			errBody, _ := ioutil.ReadAll(resp.Body)
			var errResp ErrorResponse
			json.Unmarshal(errBody, &errResp)
			log.Println(errResp.Error)
			log.Fatalln(err)
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}

		var response Response
		err = json.Unmarshal(body, &response)
		if err != nil {
			log.Fatalln(err)
		}

		// break if no more entries left
		if len(response.Content) == 0 {
			break
		}

		// loop through the validators voting powers
		for _, ele := range response.Content {
			votingPowers = append(votingPowers, ele.StakePercent)
		}

		// increment counters
		pageNo += 1
	}

	sort.Float64s(votingPowers)
	reverse(votingPowers)

	// now we're ready to calculate the nakomoto coefficient
	nakamotoCoefficient := calcNakamotoCoefficient(votingPowers)
	fmt.Println("The Nakamoto coefficient for Mina is", nakamotoCoefficient)

	return nakamotoCoefficient
}

func calcNakamotoCoefficient(votingPowers []float64) int {
	var cumulativePercent, thresholdPercent float64 = 0.00, utils.THRESHOLD_PERCENT
	nakamotoCoefficient := 0
	for _, vpp := range votingPowers {
		// since this is the  actual voting percentage, no need to multiply with 100
		cumulativePercent += vpp
		nakamotoCoefficient += 1
		if cumulativePercent >= thresholdPercent {
			break
		}
	}
	return nakamotoCoefficient
}
