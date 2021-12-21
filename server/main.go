package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
	"log"
	"os"
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

	r := gin.Default()
	coefficients := getListOfCoefficients()
	r.GET("/nakamoto-coefficients", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"coefficients": coefficients,
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func getListOfCoefficients() []int {
	var naka_coefficients []int
	var rows pgx.Rows
	var err error

	queryStmt := `SELECT naka_co_curr_val from naka_coefficients`
	if rows, err = conn.Query(context.Background(), queryStmt); err != nil {
		log.Fatalln(err)
	}
	defer rows.Close()

	for rows.Next() {
		var nc int
		err = rows.Scan(&nc)
		if err != nil {
			log.Fatalln(err)
		}
		naka_coefficients = append(naka_coefficients, nc)
	}
	return naka_coefficients
}
