package thorchain

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"sort"

	utils "github.com/xenowits/nakamoto-coefficient-calculator/utils"
)

type Response []struct {
	NodeAddress string `json:"node_address"`
	Bond        string `json:"bond"`
	Status      string `json:"status"`
}

type ErrorResponse struct {
	Id      int    `json:"id"`
	Jsonrpc string `json:"jsonrpc"`
	Error   string `json:"error"`
}

func Thorchain() (int, error) {
	votingPowers := make([]big.Int, 0, 1000)
	url := fmt.Sprintf("https://thornode.thorchain.info/thorchain/nodes")
	resp, err := http.Get(url)
	if err != nil {
		errBody, _ := ioutil.ReadAll(resp.Body)
		var errResp ErrorResponse
		json.Unmarshal(errBody, &errResp)
		log.Println(errResp.Error)
		return -1, nil
	}
	defer resp.Body.Close()

	if err != nil {
		errBody, _ := ioutil.ReadAll(resp.Body)
		var errResp ErrorResponse
		json.Unmarshal(errBody, &errResp)
		log.Println(errResp.Error)
		return -1, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return -1, err
	}

	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		return -1, err
	}

	// loop through the validators voting powers
	for _, ele := range response {
		n, ok := new(big.Int).SetString(ele.Bond, 10)
		if !ok {
			log.Fatalln("Couldn't parse string", ele.Bond)
		} else if ele.Status == "Active" {
			// Assuming we calculate only for stakers with "active" stakers
			// And discard "disabled" and "standby" stakers
			votingPowers = append(votingPowers, *n)
		}
	}

	// need to sort the powers in descending order since they are in random order
	sort.Slice(votingPowers, func(i, j int) bool {
		res := (&votingPowers[i]).Cmp(&votingPowers[j])
		if res == 1 {
			return true
		}
		return false
	})

	totalVotingPower := utils.CalculateTotalVotingPowerBigNums(votingPowers)
	fmt.Println("Total voting power:", totalVotingPower)

	// now we're ready to calculate the nakomoto coefficient
	nakamotoCoefficient := utils.CalcNakamotoCoefficientBigNums(totalVotingPower, votingPowers)
	fmt.Println("The Nakamoto coefficient for thorchain is", nakamotoCoefficient)

	return nakamotoCoefficient, nil
}
