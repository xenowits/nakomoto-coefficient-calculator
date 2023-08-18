package chains

import (
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"sort"

	utils "github.com/xenowits/nakamoto-coefficient-calculator/core/utils"
)

type Node []struct {
	Description     string `json:"description"`
	Node_Account    string `json:"node_account_id"`
	Stake 			int64  `json:"stake"`
}

type Link struct{
	Next	string `json:"next"`
}

type HederaResponse struct {
	Nodes	Node
	Links	Link
}

func Hedera() (int, error){

	// set base url for requests
	var base = "https://mainnet-public.mirrornode.hedera.com"
	var query = "/api/v1/network/nodes"

	// declare variable for tracking votes for each node
	var votingPowers []big.Int

	// declare variables for tracking pagination
	var page = ""

	//loop over api responses for all pages
	for {
		//get response from api
		resp, err := http.Get(base + query)
		if err != nil {
			fmt.Println(err)
			return 0, err
		}
		defer resp.Body.Close()

		//decode json response to go objects
		var response HederaResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		if err != nil {
			fmt.Println(err)
			return 0, err
		}
		
		// append node votes to array (from response)
		for _, node := range response.Nodes {
			votingPowers = append(votingPowers, *big.NewInt(node.Stake/100000000)) //convert tinybar to hbar
		}

		// assign next page of results to parse (null if empty, otherwise string)
		page = response.Links.Next

		// break loop where there is no more data
		if page == "" || page == "null" { break }

		//assign new query to api call and reset page variable
		query = page
		page = ""
	}

	// sort the node votes in descending order
	sort.Slice(votingPowers, func(i, j int) bool {
		return (&votingPowers[i]).Cmp(&votingPowers[j]) == 1
	})	

	// calculate the total voting power
	totalVotingPower := utils.CalculateTotalVotingPowerBigNums(votingPowers)
	fmt.Println("Total voting power for Hedera is:", new(big.Float).SetInt(totalVotingPower))

	// now we're ready to calculate the nakomoto coefficient
	nakamotoCoefficient := utils.CalcNakamotoCoefficientBigNums(totalVotingPower, votingPowers)
	fmt.Println("The Nakamoto coefficient for Hedera is", nakamotoCoefficient)

	return nakamotoCoefficient, nil	
}