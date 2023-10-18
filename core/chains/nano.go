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

type nanoOnlineRequest struct {
	Action string `json:"action"`
}

type nanoOnlineResponse struct {
	OnlineStake string `json:"trended_stake_total,omitempty"`
	Error       string `json:"error,omitempty"` // nano nodes return status 200 even when erroring
}

type nanoRepsRequest struct {
	Action  string `json:"action"`
	Sorting bool   `json:"sorting"`
	Count   int    `json:"count"`
}

type nanoRepsResponse struct {
	Reps  map[string]string `json:"representatives,omitempty"`
	Error string            `json:"error,omitempty"`
}

type nanoRichRequest struct{}

type nanoRichAccount struct {
	Balance string `json:"balance"`
	Rep     string `json:"rep"`
}

type nanoRichIdentity struct {
	Identity string                     `json:"identity"`
	Accounts map[string]nanoRichAccount `json:"accounts"`
}

type nanoRichResponse struct {
	Richlist []nanoRichIdentity `json:"richlist,omitempty"`
	Error    string             `json:"error,omitempty"`
}

type nanoEntity struct {
	Name   string
	Weight *big.Int
}

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
	rpc := "https://mynano.ninja/api/node"
	richlist := "https://nano.nendly.com/nakamoto/richlist"

	bigFromStr := func(value string) (*big.Int, error) {
		bigValue, ok := new(big.Int).SetString(value, 10)
		if !ok {
			return nil, fmt.Errorf("failed to parse big int from %v", value)
		}
		return bigValue, nil
	}

	// Determine the amount of nano that's needed for an attack
	onlineRes, err := nanoRequest[nanoOnlineRequest, nanoOnlineResponse](rpc, nanoOnlineRequest{
		Action: "confirmation_quorum",
	})

	if onlineRes != nil && onlineRes.Error != "" {
		log.Println(onlineRes.Error)
		return 0, errors.New("failed to get online stake")
	}
	if err != nil {
		return 0, err
	}

	online, err := bigFromStr(onlineRes.OnlineStake)
	if err != nil {
		return 0, err
	}

	needed := new(big.Int).Div(online, big.NewInt(3))

	// Request the list of representatives, which are what we call nano's validators (rep for short)
	repsRes, err := nanoRequest[nanoRepsRequest, nanoRepsResponse](rpc, nanoRepsRequest{
		Action:  "representatives",
		Sorting: true,
		Count:   100,
	})

	if repsRes != nil && repsRes.Error != "" {
		log.Println(repsRes.Error)
		return 0, errors.New("failed to get the list of representatives")
	}
	if err != nil {
		return 0, err
	}

	validators := make(map[string]*big.Int)
	for validator, weight := range repsRes.Reps {
		bigWeight, err := bigFromStr(weight)
		if err != nil {
			return 0, err
		}

		validators[validator] = bigWeight
	}

	// Request the "richlist" of high-value accounts
	richRes, err := nanoRequest[nanoRichRequest, nanoRichResponse](richlist, nanoRichRequest{})

	if richRes != nil && richRes.Error != "" {
		log.Println(richRes.Error)
		return 0, errors.New("failed to get the list of representatives")
	}
	if err != nil {
		return 0, err
	}

	// consider what would happen if high-value accounts moved their weights around
	for _, group := range richRes.Richlist {
		power := big.NewInt(0)

		// move all of this identity's nano off the reps
		for _, account := range group.Accounts {
			balance, err := bigFromStr(account.Balance)
			if err != nil {
				return 0, err
			}

			representative := account.Rep
			power = power.Add(power, balance)

			prior, ok := validators[representative]
			if !ok {
				continue
			}
			update := new(big.Int).Sub(prior, balance)
			if update.Sign() < 0 {
				update = big.NewInt(0)
			}
			validators[representative] = update
		}

		// appropriate other people's weight where delegated
		for account := range group.Accounts {
			if theft, ok := validators[account]; ok {
				power.Add(power, theft)
			}
			validators[account] = big.NewInt(0)
		}

		validators[group.Identity] = power
	}

	actorsByWeight := make([]nanoEntity, 0, len(validators))
	for validator, control := range validators {
		actorsByWeight = append(actorsByWeight, nanoEntity{
			Name:   validator,
			Weight: control,
		})
	}
	sort.Slice(actorsByWeight, func(i, j int) bool {
		return actorsByWeight[i].Weight.Cmp(actorsByWeight[j].Weight) > 0
	})

	// Tabulate the nakamoto coefficient
	nakamoto := 0
	colluding := big.NewInt(0)
	for _, entity := range actorsByWeight {
		nakamoto += 1
		colluding = colluding.Add(colluding, entity.Weight)
		if colluding.Cmp(needed) > 0 {
			break
		}
	}

	fmt.Println("The Nakamoto coefficient for Nano is", nakamoto)
	return nakamoto, nil
}
