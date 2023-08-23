package chains

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/big"
	"net/http"
	"sort"
	"strconv"

	utils "github.com/xenowits/nakamoto-coefficient-calculator/core/utils"
)

type SuiRespons struct {
	Jsonrpc string `json:"jsonrpc"`
	Result  struct {
		Epoch                                 string `json:"epoch"`
		ProtocolVersion                       string `json:"protocolVersion"`
		SystemStateVersion                    string `json:"systemStateVersion"`
		StorageFundTotalObjectStorageRebates  string `json:"storageFundTotalObjectStorageRebates"`
		StorageFundNonRefundableBalance       string `json:"storageFundNonRefundableBalance"`
		ReferenceGasPrice                     string `json:"referenceGasPrice"`
		SafeMode                              bool   `json:"safeMode"`
		SafeModeStorageRewards                string `json:"safeModeStorageRewards"`
		SafeModeComputationRewards            string `json:"safeModeComputationRewards"`
		SafeModeStorageRebates                string `json:"safeModeStorageRebates"`
		SafeModeNonRefundableStorageFee       string `json:"safeModeNonRefundableStorageFee"`
		EpochStartTimestampMs                 string `json:"epochStartTimestampMs"`
		EpochDurationMs                       string `json:"epochDurationMs"`
		StakeSubsidyStartEpoch                string `json:"stakeSubsidyStartEpoch"`
		MaxValidatorCount                     string `json:"maxValidatorCount"`
		MinValidatorJoiningStake              string `json:"minValidatorJoiningStake"`
		ValidatorLowStakeThreshold            string `json:"validatorLowStakeThreshold"`
		ValidatorVeryLowStakeThreshold        string `json:"validatorVeryLowStakeThreshold"`
		ValidatorLowStakeGracePeriod          string `json:"validatorLowStakeGracePeriod"`
		StakeSubsidyBalance                   string `json:"stakeSubsidyBalance"`
		StakeSubsidyDistributionCounter       string `json:"stakeSubsidyDistributionCounter"`
		StakeSubsidyCurrentDistributionAmount string `json:"stakeSubsidyCurrentDistributionAmount"`
		StakeSubsidyPeriodLength              string `json:"stakeSubsidyPeriodLength"`
		StakeSubsidyDecreaseRate              int    `json:"stakeSubsidyDecreaseRate"`
		TotalStake                            string `json:"totalStake"`
		ActiveValidators                      []struct {
			SuiAddress                   string      `json:"suiAddress"`
			ProtocolPubkeyBytes          string      `json:"protocolPubkeyBytes"`
			NetworkPubkeyBytes           string      `json:"networkPubkeyBytes"`
			WorkerPubkeyBytes            string      `json:"workerPubkeyBytes"`
			ProofOfPossessionBytes       string      `json:"proofOfPossessionBytes"`
			Name                         string      `json:"name"`
			Description                  string      `json:"description"`
			ImageURL                     string      `json:"imageUrl"`
			ProjectURL                   string      `json:"projectUrl"`
			NetAddress                   string      `json:"netAddress"`
			P2PAddress                   string      `json:"p2pAddress"`
			PrimaryAddress               string      `json:"primaryAddress"`
			WorkerAddress                string      `json:"workerAddress"`
			NextEpochProtocolPubkeyBytes interface{} `json:"nextEpochProtocolPubkeyBytes"`
			NextEpochProofOfPossession   interface{} `json:"nextEpochProofOfPossession"`
			NextEpochNetworkPubkeyBytes  interface{} `json:"nextEpochNetworkPubkeyBytes"`
			NextEpochWorkerPubkeyBytes   interface{} `json:"nextEpochWorkerPubkeyBytes"`
			NextEpochNetAddress          interface{} `json:"nextEpochNetAddress"`
			NextEpochP2PAddress          interface{} `json:"nextEpochP2pAddress"`
			NextEpochPrimaryAddress      interface{} `json:"nextEpochPrimaryAddress"`
			NextEpochWorkerAddress       interface{} `json:"nextEpochWorkerAddress"`
			VotingPower                  string      `json:"votingPower"`
			OperationCapID               string      `json:"operationCapId"`
			GasPrice                     string      `json:"gasPrice"`
			CommissionRate               string      `json:"commissionRate"`
			NextEpochStake               string      `json:"nextEpochStake"`
			NextEpochGasPrice            string      `json:"nextEpochGasPrice"`
			NextEpochCommissionRate      string      `json:"nextEpochCommissionRate"`
			StakingPoolID                string      `json:"stakingPoolId"`
			StakingPoolActivationEpoch   string      `json:"stakingPoolActivationEpoch"`
			StakingPoolDeactivationEpoch interface{} `json:"stakingPoolDeactivationEpoch"`
			StakingPoolSuiBalance        string      `json:"stakingPoolSuiBalance"`
			RewardsPool                  string      `json:"rewardsPool"`
			PoolTokenBalance             string      `json:"poolTokenBalance"`
			PendingStake                 string      `json:"pendingStake"`
			PendingTotalSuiWithdraw      string      `json:"pendingTotalSuiWithdraw"`
			PendingPoolTokenWithdraw     string      `json:"pendingPoolTokenWithdraw"`
			ExchangeRatesID              string      `json:"exchangeRatesId"`
			ExchangeRatesSize            string      `json:"exchangeRatesSize"`
		} `json:"activeValidators"`
		PendingActiveValidatorsID   string          `json:"pendingActiveValidatorsId"`
		PendingActiveValidatorsSize string          `json:"pendingActiveValidatorsSize"`
		PendingRemovals             []interface{}   `json:"pendingRemovals"`
		StakingPoolMappingsID       string          `json:"stakingPoolMappingsId"`
		StakingPoolMappingsSize     string          `json:"stakingPoolMappingsSize"`
		InactivePoolsID             string          `json:"inactivePoolsId"`
		InactivePoolsSize           string          `json:"inactivePoolsSize"`
		ValidatorCandidatesID       string          `json:"validatorCandidatesId"`
		ValidatorCandidatesSize     string          `json:"validatorCandidatesSize"`
		AtRiskValidators            []interface{}   `json:"atRiskValidators"`
		ValidatorReportRecords      [][]interface{} `json:"validatorReportRecords"`
	} `json:"result"`
	ID int `json:"id"`
}

type rawBody struct {
		JSONRPC string        `json:"jsonrpc"`
		ID      int           `json:"id"`
		Method  string        `json:"method"`
		Params  []interface{} `json:"params"`
	}

func Sui() (int, error) {
	var votingPowers []big.Int

	url := "https://fullnode.mainnet.sui.io"

	request := rawBody {
		JSONRPC: "2.0",
		ID:      1,
		Method:  "suix_getLatestSuiSystemState",
		Params:  []interface{}{},
	}

	requestBody, err := json.Marshal(request)

	if err != nil {
		log.Println(err)
	}

	// Create a new request using http
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		log.Println("Error to make request:", err)
		
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error to get response:", err)
		
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading response body:", err)
		
	}

	var response SuiRespons

	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Println("Error unmarshalling response:", err)
		
	}

	// loop through the validators voting powers, sui has voting power indicator for each validator
	// Total is 10000==100%
	for _, ele := range response.Result.ActiveValidators {
		votingPower, _ := strconv.ParseInt(ele.VotingPower, 10, 64)
		votingPowers = append(votingPowers, *big.NewInt(votingPower))
	}

	//sort
	sort.Slice(votingPowers, func(i, j int) bool {
		return votingPowers[i].Cmp(&votingPowers[j]) > 0
	})

	totalVotingPower := utils.CalculateTotalVotingPowerBigNums(votingPowers)
	fmt.Println("Total voting power for SUI:", new(big.Float).SetInt(totalVotingPower))

	// now we're ready to calculate the nakomoto coefficient
	nakamotoCoefficient := utils.CalcNakamotoCoefficientBigNums(totalVotingPower, votingPowers)
	fmt.Println("The Nakamoto coefficient for Sui is", nakamotoCoefficient)

	return nakamotoCoefficient, nil
}