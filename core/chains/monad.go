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
	"strings"
	"time"

	utils "github.com/xenowits/nakamoto-coefficient-calculator/core/utils"
)

const (
	MonadRPC     = "https://rpc.monad.xyz"
	ContractAddr = "0x0000000000000000000000000000000000001000"

	// Selectors for Monad Staking Precompile
	SelectorGetValSet  = "fb29b729"
	SelectorGetValInfo = "2b6d639a"
)

type RpcRequest struct {
	JsonRpc string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	Id      int           `json:"id"`
}

type RpcResponse struct {
	Result string `json:"result"`
	Error  struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

func Monad() (int, error) {
	// 1. Get all validator IDs via pagination
	valIDs, err := fetchAllValidatorIDs()
	if err != nil {
		return 0, err
	}
	log.Printf("Found %d active validators on Monad", len(valIDs))

	var votingPowers []big.Int

	// 2. Fetch stake for each validator
	for _, id := range valIDs {
		time.Sleep(50 * time.Millisecond)
		stake, err := fetchValidatorStake(id)
		if err != nil {
			log.Printf("Failed to fetch stake for ValID %s: %v", id.String(), err)
			continue
		}
		if stake.Cmp(big.NewInt(0)) > 0 {
			votingPowers = append(votingPowers, *stake)
		}
	}

	if len(votingPowers) == 0 {
		return 0, fmt.Errorf("no voting power found after querying %d validators", len(valIDs))
	}

	// Sort by stake descending
	sort.Slice(votingPowers, func(i, j int) bool {
		return (&votingPowers[i]).Cmp(&votingPowers[j]) > 0
	})

	totalStake := utils.CalculateTotalVotingPowerBigNums(votingPowers)

	fmt.Println("Total Monad Stake:", new(big.Float).SetInt(totalStake))

	nc := utils.CalcNakamotoCoefficientBigNums(totalStake, votingPowers)
	fmt.Println("Monad Nakamoto Coefficient:", nc)

	return nc, nil
}

// fetchAllValidatorIDs paginates through the system contract to retrieve all validator IDs.
func fetchAllValidatorIDs() ([]*big.Int, error) {
	var allIDs []*big.Int
	currentIndex := 0

	for {
		// Payload: selector + current_index (uint256 encoded)
		arg := fmt.Sprintf("%064x", currentIndex)
		data := "0x" + SelectorGetValSet + arg

		res, err := ethCall(data)
		if err != nil {
			return nil, err
		}

		res = strings.TrimPrefix(res, "0x")

		// Response ABI: [bool is_done, uint256 next_index, uint256 offset, uint256 length, ...items]
		if len(res) < 256 {
			return nil, fmt.Errorf("response too short during pagination")
		}

		// Parse is_done (Word 0)
		isDoneHex := res[0:64]
		isDoneVal, _ := new(big.Int).SetString(isDoneHex, 16)
		isDone := isDoneVal.Cmp(big.NewInt(1)) == 0

		// Parse next_index (Word 1)
		nextIndexHex := res[64:128]
		nextIndexBig, _ := new(big.Int).SetString(nextIndexHex, 16)
		currentIndex = int(nextIndexBig.Int64())

		// Parse array length (Word 3)
		lenHex := res[192:256]
		countBig, _ := new(big.Int).SetString(lenHex, 16)
		count := int(countBig.Int64())

		// Extract items starting at offset 256
		dataStart := 256
		for i := 0; i < count; i++ {
			p := dataStart + (i * 64)
			if p+64 > len(res) {
				break
			}
			segment := res[p : p+64]
			val := new(big.Int)
			val.SetString(segment, 16)
			allIDs = append(allIDs, val)
		}

		if isDone {
			break
		}
	}

	return allIDs, nil
}

func fetchValidatorStake(valID *big.Int) (*big.Int, error) {
	// Payload: selector + val_id (uint256 encoded)
	arg := fmt.Sprintf("%064x", valID)
	data := "0x" + SelectorGetValInfo + arg

	res, err := ethCall(data)
	if err != nil {
		return nil, err
	}

	res = strings.TrimPrefix(res, "0x")

	// ValidatorInfo struct layout: Stake is at index 6 (offset 192 bytes / 384 hex chars)
	start := 384
	if len(res) < start+64 {
		return nil, fmt.Errorf("response too short for validator info")
	}

	stakeHex := res[start : start+64]
	stake := new(big.Int)
	stake.SetString(stakeHex, 16)

	return stake, nil
}

func ethCall(data string) (string, error) {
	payload := RpcRequest{
		JsonRpc: "2.0",
		Method:  "eth_call",
		Params: []interface{}{
			map[string]string{
				"to":   ContractAddr,
				"data": data,
			},
			"latest",
		},
		Id: 1,
	}

	bodyBytes, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", MonadRPC, bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Nakaflow/1.0")

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	if bytes.HasPrefix(bytes.TrimSpace(respBody), []byte("<")) {
		return "", fmt.Errorf("RPC returned HTML content")
	}

	var rpcResp RpcResponse
	if err := json.Unmarshal(respBody, &rpcResp); err != nil {
		return "", fmt.Errorf("json parse error: %v", err)
	}

	if rpcResp.Error.Code != 0 {
		return "", fmt.Errorf("rpc error: %s", rpcResp.Error.Message)
	}

	return rpcResp.Result, nil
}
