package main

import (
	"fmt"
	"log"
	"sort"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xenowits/nakamoto-coefficient-calculator/core/chains"
)

type JsonResponse struct {
	ChainName     string `json:"chain_name"`
	ChainToken    string `json:"chain_token"`
	NakaCoPrevVal int    `json:"naka_co_prev_val"`
	NakaCoCurrVal int    `json:"naka_co_curr_val"`
	Change        int    `json:"naka_co_change_val"`
}

func main() {
	var mu sync.Mutex
	chainState := chains.NewState()

	// Run a goroutine which refreshes state after every interval.
	ticker := time.NewTicker(6 * time.Hour)
	quit := make(chan struct{})
	defer close(quit)

	go func(state chains.ChainState) {
		for {
			select {
			case <-ticker.C:
				log.Println("Ticker ticked")
				newState := chains.RefreshChainState(chainState)

				mu.Lock()
				chainState = newState
				mu.Unlock()

				fmt.Println(chainState)
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}(chainState)

	// Run server.
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.GET("/naka-coeffs", func(c *gin.Context) {
		coefficients := getListOfCoefficients(chainState)
		c.Header("Access-Control-Allow-Origin", "*")
		c.JSON(200, gin.H{
			"coefficients": coefficients,
		})
	})
	r.Run(":8080") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func getListOfCoefficients(state chains.ChainState) []JsonResponse {
	var coeffs []JsonResponse
	for token, chain := range state {
		coeffs = append(coeffs, JsonResponse{
			ChainName:     token.ChainName(),
			ChainToken:    string(token),
			NakaCoPrevVal: chain.PrevNCVal,
			NakaCoCurrVal: chain.CurrNCVal,
			Change:        chain.CurrNCVal - chain.PrevNCVal,
		})
	}

	sort.Slice(coeffs, func(i, j int) bool {
		if coeffs[i].ChainToken < coeffs[j].ChainToken {
			return true
		}

		return false
	})

	return coeffs
}
