package chains

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"

	"github.com/xenowits/nakamoto-coefficient-calculator/core/utils"
)

const totalValidatorsUrl = "https://api.multiversx.com/stake"
const identitiesUrl = "https://api.multiversx.com/identities"

type MultiversXTotalValidatorsResponse struct {
	TotalValidators int64 `json:"totalValidators"`
}

type MultiversXIdentitiesResponse []struct {
	Locked        string `json:"locked"`
	NumValidators int64  `json:"validators"`
}

func MultiversX() (int, error) {
	numValidatorsPerIdentity := make([]int64, 0)

	totalNumberOfValidators, err := getTotalValidatorsNumber()
	if err != nil {
		return 0, err
	}

	identities, err := getIdentities()
	if err != nil {
		return 0, err
	}

	for _, identity := range identities {
		if identity.Locked == "0" {
			continue
		}
		numValidatorsPerIdentity = append(numValidatorsPerIdentity, identity.NumValidators)
	}

	sort.Slice(numValidatorsPerIdentity, func(i, j int) bool {
		return numValidatorsPerIdentity[i] > numValidatorsPerIdentity[j]
	})

	fmt.Println("Total voting power:", totalNumberOfValidators)

	// there is a fixed number of validator seats in MultiversX - currently 3200
	// the Nakamoto coefficient can be computed by counting the identities (node operators)
	// that control more than 33% of the total number of validators
	nakamotoCoefficient := utils.CalcNakamotoCoefficient(totalNumberOfValidators, numValidatorsPerIdentity)
	fmt.Println("The Nakamoto coefficient for MultiversX is", nakamotoCoefficient)

	return nakamotoCoefficient, nil
}

func getTotalValidatorsNumber() (int64, error) {
	resp, err := http.Get(totalValidatorsUrl)
	if err != nil {
		return 0, err
	}

	defer closeBody(resp)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	var response MultiversXTotalValidatorsResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return 0, err
	}

	return response.TotalValidators, nil
}

func getIdentities() (MultiversXIdentitiesResponse, error) {
	resp, err := http.Get(identitiesUrl)
	if err != nil {
		return nil, err
	}

	defer closeBody(resp)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response MultiversXIdentitiesResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func closeBody(resp *http.Response) {
	if resp == nil || resp.Body == nil {
		return
	}
	if closeErr := resp.Body.Close(); closeErr != nil {
		log.Printf("failed to close response body: %s", closeErr)
	}
}
