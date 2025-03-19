package chains

import (
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"sort"
	"strconv"

	"github.com/xenowits/nakamoto-coefficient-calculator/core/utils"
)

type NanExplorerResponse struct {
	Rep []struct {
		Account string `json:"account"`
		Weight  string `json:"weight"`
	} `json:"rep"`
	OnlineStakeTotal string `json:"online_stake_total"`
}

type EntityResponse struct {
	Timestamp int64 `json:"timestamp"`
	Entities  []struct {
		Entity          string   `json:"entity"`
		Representatives []string `json:"representatives"`
	} `json:"entities"`
}

const (
	THRESHOLD = 67 // 67% threshold for Nakamoto Coefficient
)

func Nano() (int, error) {

	// Step 1: Fetch entity groups
	resp, err := http.Get("https://nanocharts.info/data/entities.json")
	if err != nil {
		log.Println("Error fetching entities:", err)
		return 0, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Printf("Entities fetch failed: %d", resp.StatusCode)
		return 0, fmt.Errorf("entities fetch failed: %d", resp.StatusCode)
	}

	var entityData EntityResponse
	if err := json.NewDecoder(resp.Body).Decode(&entityData); err != nil {
		log.Println("Error decoding entities JSON:", err)
		return 0, err
	}

	entityGroups := make(map[string][]string)
	for _, entity := range entityData.Entities {
		entityGroups[entity.Entity] = entity.Representatives
	}

	// Step 2: Fetch online reps and weights from NanExplorer
	resp, err = http.Get("https://api.nanexplorer.com/representatives_online?network=nano")
	if err != nil {
		log.Println("Error fetching online reps:", err)
		return 0, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Printf("NanExplorer fetch failed: %d", resp.StatusCode)
		return 0, fmt.Errorf("nanexplorer fetch failed: %d", resp.StatusCode)
	}

	var explorerData NanExplorerResponse
	if err := json.NewDecoder(resp.Body).Decode(&explorerData); err != nil {
		log.Println("Error decoding nanexplorer data:", err)
		return 0, err
	}

	// Step 3: Process weights into big.Int
	weights := make(map[string]*big.Int)
	accountToEntity := make(map[string]string)

	for entity, reps := range entityGroups {
		for _, acc := range reps {
			accountToEntity[acc] = entity
		}
	}

	for _, rep := range explorerData.Rep {
		weight, err := strconv.ParseFloat(rep.Weight, 64)
		if err != nil {
			log.Printf("Error parsing weight for %s: %v", rep.Account, err)
			continue
		}
		weightInt := new(big.Int).SetInt64(int64(weight * 1e6)) // Convert XNO to raw-like integer

		entityName, ok := accountToEntity[rep.Account]
		if !ok {
			entityName = rep.Account
		}
		if _, ok := weights[entityName]; !ok {
			weights[entityName] = new(big.Int)
		}
		weights[entityName].Add(weights[entityName], weightInt)
	}

	// Step 4: Collect voting powers as big.Int (not pointers)
	var votingPowers []big.Int
	for _, weight := range weights {
		// Copy the value of weight to the slice of big.Int
		votingPowers = append(votingPowers, *weight)
	}

	if len(votingPowers) == 0 {
		log.Println("No weights processed - no online reps")
		return 0, fmt.Errorf("no weights")
	}

	// Manually accumulate voting power until we hit the threshold

	// Calculate total voting power
	calculatedTotalVotingPower := utils.CalculateTotalVotingPowerBigNums(votingPowers)
	thresholdVotingPower := new(big.Int).Mul(calculatedTotalVotingPower, big.NewInt(THRESHOLD))
	thresholdVotingPower.Div(thresholdVotingPower, big.NewInt(100))

	// Sort the voting powers in descending order
	sort.Slice(votingPowers, func(i, j int) bool {
		return votingPowers[i].Cmp(&votingPowers[j]) > 0
	})

	// Step 5: Accumulate until the threshold is met
	var accumulatedVotingPower big.Int
	for i, power := range votingPowers {
		accumulatedVotingPower.Add(&accumulatedVotingPower, &power)
		if accumulatedVotingPower.Cmp(thresholdVotingPower) >= 0 {

			log.Printf("Nakamoto Coefficient (67%%): %d", i+1) // Number of entities needed to meet threshold

			return i + 1, nil
		}
	}

	// In case we have all entities
	// log.Printf("Nakamoto Coefficient (67%%): %d", len(votingPowers))
	return len(votingPowers), nil
}
