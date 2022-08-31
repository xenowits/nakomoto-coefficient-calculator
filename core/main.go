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
		fmt.Fprintln(os.Stderr, "Unable to connect to database", err)
		os.Exit(1)
	}

	defer func() {
		if err := conn.Close(context.Background()); err != nil {
			log.Println("failed to close database connection", err)
		}
	}()

	networks := []string{"BLD", "REGEN", "ETH", "BNB", "ATOM", "OSMO", "MATIC", "MINA", "SOL", "AVAX", "LUNA", "GRT", "RUNE", "NEAR", "JUNO", "NANO"}
	for _, n := range networks {
		UpdateChainInfo(n)
	}
}

func UpdateChainInfo(chainToken string) {
	prevVal, currVal := getPrevVal(chainToken), 0
	var err error
	switch chainToken {
	case "BLD":
		currVal, err = chains.Agoric()
	case "ETH2":
		currVal, err = chains.Eth()
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
	case "NEAR":
		currVal, err = chains.Near()
	case "JUNO":
		currVal, err = chains.Juno()
	case "REGEN":
		currVal, err = chains.Regen()
	case "NANO":
		currVal, err = chains.Nano()
	}

	if err != nil {
		log.Println("failed to update chain info", chainToken, err)
	}

	if err := saveUpdatedVals(currVal, prevVal, chainToken); err != nil {
		log.Println("failed to save updated values to database", chainToken, err)
	}
}

// GetPrevVal queries the database to get the previous (prior to updating it now) value of nakamoto coefficient for the given chain
// Assumes row for chain already exists in the table
func getPrevVal(chainToken string) int {
	queryStmt := `SELECT naka_co_curr_val from naka_coefficients WHERE chain_token=$1`
	var nakaCoeffPrevVal int
	if err := conn.QueryRow(context.Background(), queryStmt, chainToken).Scan(&nakaCoeffPrevVal); err == nil {
	} else {
		fmt.Println("Read unsuccessful", chainToken, err)
		return -1
	}
	return nakaCoeffPrevVal
}

// SaveUpdatedVals saves the recently calculated values back to the database
func saveUpdatedVals(currVal int, prevVal int, chainToken string) error {
	queryStmt := `UPDATE naka_coefficients SET naka_co_curr_val=$1, naka_co_prev_val=$2 WHERE chain_token=$3`
	_, err := conn.Exec(context.Background(), queryStmt, currVal, prevVal, chainToken)
	if err != nil {
		fmt.Println("Write unsuccessful for "+chainToken, err)
	}
	return err
}
