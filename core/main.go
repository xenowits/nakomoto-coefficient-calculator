package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v4"
	"github.com/xenowits/nakamoto-coefficient-calculator/core/chains"
)

var conn *pgx.Conn

func main() {
	var err error
	conn, err = pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	networks := []string{"BNB", "ATOM", "OSMO", "MATIC", "MINA", "SOL", "AVAX", "LUNA", "GRT", "RUNE"}

	for _, n := range networks {
		UpdateChainInfo(n)
	}
}

func UpdateChainInfo(chain_token string) {
	prevVal, currVal := GetPrevVal(chain_token), 0
	var err error
	switch chain_token {
	case "BNB":
		currVal, err = chains.Binance()
	case "ATOM":
		currVal, err = chains.Cosmos()
	case "OSMO":
		currVal, err = chains.Osmosis()
	case "MATIC":
		currVal, err = chains.Polygon()
	case "MINA":
		currVal, err = chains.Mina()
	case "SOL":
		currVal, err = chains.Solana()
	case "AVAX":
		currVal, err = chains.Avalanche()
	case "LUNA":
		currVal, err = chains.Terra()
	case "GRT":
		currVal, err = chains.Graph()
	case "RUNE":
		currVal, err = chains.Thorchain()
	}

	if err != nil {
		log.Println("Error occurred for", chain_token, err)
	} else {
		SaveUpdatedVals(currVal, prevVal, chain_token)
	}
}

// Query the database to get the previous (prior to updating it now) value of nakamoto coefficient for the given chain
func GetPrevVal(chain_token string) int {
	queryStmt := `SELECT naka_co_curr_val from naka_coefficients WHERE chain_token=$1`
	var naka_co_prev_val int
	if err := conn.QueryRow(context.Background(), queryStmt, chain_token).Scan(&naka_co_prev_val); err == nil {
	} else {
		fmt.Println("Read unsuccessful for "+chain_token, err)
		return -1
	}
	return naka_co_prev_val
}

// Save the recently calculated values back to the database
func SaveUpdatedVals(curr_val int, prev_val int, chain_token string) error {
	queryStmt := `UPDATE naka_coefficients SET naka_co_curr_val=$1, naka_co_prev_val=$2 WHERE chain_token=$3`
	_, err := conn.Exec(context.Background(), queryStmt, curr_val, prev_val, chain_token)
	if err != nil {
		fmt.Println("Write unsuccessful for "+chain_token, err)
	}
	return err
}
