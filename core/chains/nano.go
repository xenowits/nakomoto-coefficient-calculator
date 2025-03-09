package chains

import (
	"encoding/json"
	"log"
	"net/http"
)

// NanoStats represents the structure for Nano's stats data from nanonakamoto.xyz.
type NanoStats struct {
     NakamotoCoefficient int `json:"nakamotoCoefficient"`
}

func Nano() (int, error) {
	chartsURL := "https://nanonakamoto.xyz/api/nc"

	log.Println("Fetching data for Nano")

	// Fetch the data from nanocharts
	resp, err := http.Get(chartsURL)
	if err != nil {
		log.Println("Error fetching data from nanonakamoto.xyz:", err)
		return 0, err
	}
	defer resp.Body.Close()

	var data NanoStats

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		log.Println("Error decoding JSON:", err)
		return 0, err
	}

	// The Nakamoto coefficient is directly available in the JSON!
	nakamotoCoefficient := data.NakamotoCoefficient

	log.Println("The Nakamoto coefficient for Nano is", nakamotoCoefficient)

	return nakamotoCoefficient, nil
}
