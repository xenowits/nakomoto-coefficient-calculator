package chains

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const nakamotoThreshold = 33

type celestiaResp struct {
	Jailed             bool    `json:"jailed"`
	VotingPowerPercent float64 `json:"votingPowerPercent"`
}

func Celestia() (int, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFunc()

	url := "https://celestia.api.explorers.guru/api/v1/validators"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)

	if err != nil {
		return 0, err
	}

	resp, err := new(http.Client).Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	var response []celestiaResp
	err = json.Unmarshal(body, &response)
	if err != nil {
		return 0, err
	}

	var (
		cumulativePower     float64
		nakamotoCoefficient int
	)
	for _, resp := range response {
		cumulativePower += resp.VotingPowerPercent
		nakamotoCoefficient += 1
		if cumulativePower > nakamotoThreshold {
			break
		}
	}

	fmt.Printf("The Nakamoto coefficient for %s is %d\n", "celestia", nakamotoCoefficient)

	return nakamotoCoefficient, nil
}
