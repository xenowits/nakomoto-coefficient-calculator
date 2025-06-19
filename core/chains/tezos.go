package chains

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	tzktTotalBakingPowerURL = "https://api.tzkt.io/v1/cycles?select=index,totalBakingPower&limit=1"
	tzktDelegatesURL    = "https://api.tzkt.io/v1/delegates?active=true&select=address,alias,totalStakedBalance,delegatedBalance&sort.desc=stakingBalance&limit=100&offset=%d"
)

type TezosCycle struct {
	Index        int   `json:"index"`
	TotalBakingPower int64 `json:"totalBakingPower"`
}

type TezosBaker struct {
	Address        string `json:"address"`
	Alias          string `json:"alias"`
	StakingBalance int64  `json:"totalStakedBalance"`
	DelegatedBalance int64 `json:"delegatedBalance"`
}

func Tezos() (int, error) {
	// Step 1: Fetch total baking power for current cycle
	resp, err := http.Get(tzktTotalBakingPowerURL)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var cycles []TezosCycle
	if err := json.NewDecoder(resp.Body).Decode(&cycles); err != nil || len(cycles) == 0 {
		return 0, fmt.Errorf("failed to parse total baking power")
	}
	TotalBakingPower := cycles[0].TotalBakingPower
	threshold := float64(TotalBakingPower) * 0.33

	var (
		offset      = 0
		limit       = 10
		accumulated int64
		coefficient int
	)

	// Step 2: Accumulate baking power for top bakers
	for accumulated < int64(threshold) {
		url := fmt.Sprintf(tzktDelegatesURL, offset)
		resp, err := http.Get(url)
		if err != nil {
			return 0, err
		}
		defer resp.Body.Close()

		var bakers []TezosBaker
		if err := json.NewDecoder(resp.Body).Decode(&bakers); err != nil {
			return 0, err
		}
		if len(bakers) == 0 {
			break
		}

		for _, baker := range bakers {
			var bakerPower = baker.StakingBalance + (baker.DelegatedBalance / 3)
			accumulated += bakerPower
			coefficient++

			if float64(accumulated) >= threshold {
				return coefficient, nil
			}
		}

		offset += limit
	}

	return coefficient, nil
}
