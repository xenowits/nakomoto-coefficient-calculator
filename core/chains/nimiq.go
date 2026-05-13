package chains

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const nimiqNakamotoEndpoint = "https://nimiq.fyi/api/nakamoto-coefficient/nimiq"

type nimiqNakamotoResponse struct {
	Value        int    `json:"value"`
	CurrentValue int    `json:"naka_co_curr_val"`
	ChainName    string `json:"chain_name"`
	ChainToken   string `json:"chain_token"`
	CalculatedAt string `json:"calculated_at"`
}

var nimiqClient = &http.Client{Timeout: 10 * time.Second}

func Nimiq() (int, error) {
	req, err := http.NewRequest("GET", nimiqNakamotoEndpoint, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to create Nimiq request: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "Nakaflow/nakamoto-coefficient-calculator")

	resp, err := nimiqClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch Nimiq Nakamoto coefficient: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("Nimiq API returned status: %d", resp.StatusCode)
	}

	var payload nimiqNakamotoResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return 0, fmt.Errorf("failed to decode Nimiq response: %w", err)
	}

	value := payload.Value
	if value == 0 {
		value = payload.CurrentValue
	}
	if value <= 0 {
		return 0, fmt.Errorf("Nimiq API returned invalid Nakamoto coefficient: %d", value)
	}
	return value, nil
}
