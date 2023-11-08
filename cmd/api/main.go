package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Osagie-Godstand/chi-postgres-user-account/db"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	config := &db.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Password: os.Getenv("DB_PASSWORD"),
		User:     os.Getenv("DB_USER"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
		DBName:   os.Getenv("DB_NAME"),
	}

	dbConn, err := db.NewConnection(config)
	if err != nil {
		log.Fatal("could not connect to the database:", err)
	}

	migrationsErr := db.RunMigrations(dbConn)
	if migrationsErr != nil {
		log.Fatal("could not migrate the database:", migrationsErr)
	}

	router := initializeRouter(dbConn)
	listenAddr := os.Getenv("HTTP_LISTEN_ADDRESS")

	// Starting the HTTP server
	log.Printf("Server is listening on %s...", listenAddr)
	if err := http.ListenAndServe(listenAddr, router); err != nil {
		log.Fatal("HTTP server error:", err)
	}
}
