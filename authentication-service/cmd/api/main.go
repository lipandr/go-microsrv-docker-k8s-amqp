package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"authentication/data"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const webPort = "80"

var counts int

type Config struct {
	DB     *sql.DB
	Models data.Models
}

func main() {
	fmt.Println("Starting authentication service")

	// connect to database
	conn := connectToDB()
	if conn == nil {
		log.Panic("Can't connect to database")
	}
	app := Config{
		DB:     conn,
		Models: data.New(conn),
	}

	srv := http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Panic(err)
	}
}

func openDB(dsn string) (*sql.DB, error) {
	// Open a connection to the database
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	// Ping the connection to make sure it's alive
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func connectToDB() *sql.DB {
	dsn := os.Getenv("DSN")

	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Println("Postgres is not yet ready ...")
			counts++
		} else {
			log.Println("Connected to Postgres!")
			return connection
		}
		if counts > 10 {
			log.Println(err)
		}
		log.Println("Backing off for two seconds ...")
		time.Sleep(2 * time.Second)
	}
}
