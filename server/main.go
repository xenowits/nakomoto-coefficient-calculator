package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
)

var conn *pgx.Conn

type JsonResponse struct {
	Chain_name       string `json:"chain_name"`
	Chain_token      string `json:"chain_token"`
	Naka_co_curr_val int    `json:"naka_co_curr_val"`
	Naka_co_prev_val int    `json:"naka_co_prev_val"`
	Change           int    `json:"naka_co_change_val"`
}

func main() {
	var err error
	conn, err = pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.GET("/nakamoto-coefficients", func(c *gin.Context) {
		coefficients := getListOfCoefficients()
		c.Header("Access-Control-Allow-Origin", "*")
		c.JSON(200, gin.H{
			"coefficients": coefficients,
		})
	})
	r.Run(":8080") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func getListOfCoefficients() []JsonResponse {
	var naka_coefficients []JsonResponse
	var rows pgx.Rows
	var err error

	queryStmt := `SELECT chain_name, chain_token, naka_co_prev_val, naka_co_curr_val from naka_coefficients`
	if rows, err = conn.Query(context.Background(), queryStmt); err != nil {
		log.Fatalln(err)
	}
	defer rows.Close()

	for rows.Next() {
		var chain_name, chain_token string
		var nc_prev_val, nc_curr_val int
		err = rows.Scan(&chain_name, &chain_token, &nc_prev_val, &nc_curr_val)
		if err != nil {
			log.Fatalln(err)
		}
		naka_coefficients = append(naka_coefficients, JsonResponse{chain_name, chain_token, nc_curr_val, nc_prev_val, nc_curr_val - nc_prev_val})
	}
	return naka_coefficients
}
