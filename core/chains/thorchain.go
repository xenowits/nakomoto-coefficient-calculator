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

	"github.com/xenowits/nakamoto-coefficient-calculator/core/utils"
)

type ThorchainResponse []struct {
	NodeAddress string `json:"node_address"`
	Bond        string `json:"total_bond"`
	Status      string `json:"status"`
}

type ThorchainErrorResponse struct {
	Id      int    `json:"id"`
	Jsonrpc string `json:"jsonrpc"`
	Error   string `json:"error"`
}

func Thorchain() (int, error) {
	votingPowers := make([]big.Int, 0, 1000)
	url := fmt.Sprintf("https://thornode.ninerealms.com/thorchain/nodes")
	ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFunc()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.Println(err)
		return 0, errors.New("create get request for thorchain")
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

	var response ThorchainResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return 0, err
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
