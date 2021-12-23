package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4"
	"github.com/xenowits/nakamoto-coefficient-calculator/cmd/binance"
	"github.com/xenowits/nakamoto-coefficient-calculator/cmd/cosmos"
	// "github.com/xenowits/nakamoto-coefficient-calculator/cmd/mina"
	"github.com/xenowits/nakamoto-coefficient-calculator/cmd/osmosis"
	"github.com/xenowits/nakamoto-coefficient-calculator/cmd/polygon"
	"github.com/xenowits/nakamoto-coefficient-calculator/cmd/solana"
	"github.com/xenowits/nakamoto-coefficient-calculator/cmd/avalanche"
)

var conn *pgx.Conn

func main() {
	var err error
	// Note: Uncomment the following line & set the database_url secretly, please.
	os.Setenv("DATABASE_URL", "postgres://xenowits:xenowits@localhost:5432/postgres")
	conn, err = pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	// binance
	prevVal := GetPrevVal("BNB")
	currVal := binance.Binance()
	saveUpdatedVals(currVal, prevVal, "BNB")

	// cosmos
	prevVal = GetPrevVal("ATOM")
	currVal = cosmos.Cosmos()
	saveUpdatedVals(currVal, prevVal, "ATOM")

	// osmosis
	prevVal = GetPrevVal("OSMO")
	currVal = osmosis.Osmosis()
	saveUpdatedVals(currVal, prevVal, "OSMO")

	// polygon
	prevVal = GetPrevVal("MATIC")
	currVal = polygon.Polygon()
	saveUpdatedVals(currVal, prevVal, "MATIC")

	// mina
	// prevVal = GetPrevVal("MINA")
	// currVal = mina.Mina()
	// saveUpdatedVals(currVal, prevVal, "MINA")

	// solana
	prevVal = GetPrevVal("SOL")
	currVal = solana.Solana()
	saveUpdatedVals(currVal, prevVal, "SOL")

	// solana
	prevVal = GetPrevVal("AVAX")
	currVal = avalanche.Avalanche()
	saveUpdatedVals(currVal, prevVal, "AVAX")
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
func saveUpdatedVals(curr_val int, prev_val int, chain_token string) error {
	queryStmt := `UPDATE naka_coefficients SET naka_co_curr_val=$1, naka_co_prev_val=$2 WHERE chain_token=$3`
	_, err := conn.Exec(context.Background(), queryStmt, curr_val, prev_val, chain_token)
	if err != nil {
		fmt.Println("Write unsuccessful for "+chain_token, err)
	}
	return err
}
