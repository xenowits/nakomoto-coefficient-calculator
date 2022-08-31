package chains

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math/big"
	"net/http"
	"sort"
)

func nanoRequest[ReqType any, ResType any](url string, req ReqType) (*ResType, error) {
	var response ResType
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(body, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

func Nano() (int, error) {
	rpc := "http://localhost:7076"
	richlist := "https://nano.nendly.com/nakamoto/richlist"

	bigFromStr := func(value string) *big.Int {
		bigValue, _ := new(big.Int).SetString(value, 10)
		return bigValue
	}

	// Determine the amount of nano that's needed for an attack
	type OnlineRequest struct {
		Action string `json:"action"`
	}
	type OnlineResponse struct {
		OnlineStake string `json:"trended_stake_total,omitempty"`
		Error       string `json:"error,omitempty"` // nano nodes return status 200 even when erroring
	}
	onlineRes, err := nanoRequest[OnlineRequest, OnlineResponse](rpc, OnlineRequest{
		Action: "confirmation_quorum",
	})
	if onlineRes != nil && onlineRes.Error != "" {
		log.Println(onlineRes.Error)
		return -1, errors.New("failed to get the online stake")
	}
	if err != nil {
		return -1, err
	}
	online := bigFromStr(onlineRes.OnlineStake)
	needed := new(big.Int).Div(online, big.NewInt(3))

	// Request the list of reps
	type RepsRequest struct {
		Action  string `json:"action"`
		Sorting bool   `json:"sorting"`
		Count   int    `json:"count"`
	}
	type RepsResponse struct {
		Reps  map[string]string `json:"representatives,omitempty"`
		Error string            `json:"error,omitempty"`
	}
	repsRes, err := nanoRequest[RepsRequest, RepsResponse](rpc, RepsRequest{
		Action:  "representatives",
		Sorting: true,
		Count:   100,
	})
	if repsRes != nil && repsRes.Error != "" {
		log.Println(repsRes.Error)
		return -1, errors.New("failed to get the list of representatives")
	}
	if err != nil {
		return -1, err
	}
	reps := make(map[string]*big.Int)
	for rep, weight := range repsRes.Reps {
		bigWeight := bigFromStr(weight)
		reps[rep] = bigWeight
	}

	// Request the "richlist" of high-value accounts
	type RichRequest struct{}
	type RichAccount struct {
		Balance string `json:"balance"`
		Rep     string `json:"rep"`
	}
	type RichIdentity struct {
		Identity string                 `json:"identity"`
		Accounts map[string]RichAccount `json:"accounts"`
	}
	type RichResponse struct {
		Richlist []RichIdentity `json:"richlist,omitempty"`
		Error    string         `json:"error,omitempty"`
	}
	richRes, err := nanoRequest[RichRequest, RichResponse](richlist, RichRequest{})
	if richRes != nil && richRes.Error != "" {
		log.Println(richRes.Error)
		return -1, errors.New("failed to get the list of representatives")
	}
	if err != nil {
		return -1, err
	}

	// consider what would happen if high-value accounts moved their weights around
	for _, group := range richRes.Richlist {
		power := big.NewInt(0)

		// move all of this identity's nano off the reps
		for _, account := range group.Accounts {
			balance := bigFromStr(account.Balance)
			rep := account.Rep
			power = power.Add(power, balance)

			prior, ok := reps[rep]
			if !ok {
				continue
			}
			update := new(big.Int).Sub(prior, balance)
			if update.Sign() < 0 {
				update = big.NewInt(0)
			}
			reps[rep] = update
		}

		// appropriate other people's weight where delegated
		for account := range group.Accounts {
			if theft, ok := reps[account]; ok {
				power.Add(power, theft)
			}
			reps[account] = big.NewInt(0)
		}

		reps[group.Identity] = power
	}

	type Entity struct {
		Name   string
		Weight *big.Int
	}

	topReps := make([]Entity, 0, len(reps))
	for rep, control := range reps {
		topReps = append(topReps, Entity{
			Name:   rep,
			Weight: control,
		})
	}
	sort.Slice(topReps, func(i, j int) bool {
		return topReps[i].Weight.Cmp(topReps[j].Weight) > 0
	})

	// Tabulate the nakamoto coefficient
	nakamoto := 0
	colluding := big.NewInt(0)
	for _, entity := range topReps {
		nakamoto += 1
		colluding = colluding.Add(colluding, entity.Weight)
		if colluding.Cmp(needed) > 0 {
			break
		}
	}

	fmt.Println("The Nakamoto coefficient for Nano is", nakamoto)
	return nakamoto, nil
}
