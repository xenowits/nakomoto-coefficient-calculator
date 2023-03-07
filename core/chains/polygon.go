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

type PolygonResponse struct {
	List []struct {
		TotalStaked int64 `json:"totalStaked"`
	} `json:"list"`
}

func Polygon() (int, error) {
	var votingPowers []int64
	ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFunc()

	url := fmt.Sprintf("https://validator.info/api/polygon/validators?timeframe=week&nameContains=&activeValidators=true")
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.Println(err)
		return 0, errors.New("create get request for polygon")
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

	var response PolygonResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return 0, err
	}

	// Loop through the validators staked amounts
	for _, ele := range response.List {
		votingPowers = append(votingPowers, ele.TotalStaked)
	}

	// Sort the voting powers in descending order since they maybe in random order.
	sort.Slice(votingPowers, func(i, j int) bool { return votingPowers[i] > votingPowers[j] })

	totalVotingPower := utils.CalculateTotalVotingPower(votingPowers)
	fmt.Println("Total voting power:", totalVotingPower)

	// Now we're ready to calculate the nakamoto coefficient
	nakamotoCoefficient := utils.CalcNakamotoCoefficient(totalVotingPower, votingPowers)
	fmt.Println("The Nakamoto coefficient for 0xPolygon is", nakamotoCoefficient)

	return nakamotoCoefficient, nil
}
