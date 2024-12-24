package chains

import (
	"encoding/json"
	"log"
	"net/http"
)

func Nano() (int, error) {
	chartsURL := "https://nanocharts.info/data/nanocharts.json"

	log.Println("Fetching data for Nano")

	// Fetch the data from nanocharts
	resp, err := http.Get(chartsURL)
	if err != nil {
		log.Println("Error fetching data from nanocharts:", err)
		return 0, err
	}
	defer resp.Body.Close()

	var data struct {
		Stats struct {
			C1s struct {
				N int `json:"n"`
			} `json:"c1s"`
		} `json:"stats"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		log.Println("Error decoding JSON:", err)
		return 0, err
	}

	// The Nakamoto coefficient is directly available in the JSON under stats.c1s.n
	nakamotoCoefficient := data.Stats.C1s.N

	log.Println("The Nakamoto coefficient for Nano is", nakamotoCoefficient)

	return nakamotoCoefficient, nil
}
