package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/Osagie-Godstand/chi-user-account-crelog/db"
	"github.com/Osagie-Godstand/chi-user-account-crelog/handlers"
	"github.com/go-chi/chi/v5"
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

func initializeRouter(dbConn *sql.DB) http.Handler {
	router := chi.NewRouter()

	// Initializing the repositories and handlers
	userRepository := db.NewPostgresUserRepository(dbConn)
	session := handlers.NewSessionHandler(dbConn, userRepository)
	userHandler := handlers.NewUserHandler(userRepository)

	// Defining routes and handlers
	router.Post("/login", session.Login)
	router.Post("/user", userHandler.HandlePostUser)

	// Applying the ValidateSession middleware to routes that need session validation
	router.With(session.ValidateSession).Post("/logout", session.Logout)
	router.With(session.ValidateSession).Put("/user/{id}", userHandler.HandlePutUser)
	router.With(session.ValidateSession).Get("/users", userHandler.HandleGetUsers)
	router.With(session.ValidateSession).Get("/user/{id}", userHandler.HandleGetUserByID)
	router.With(session.ValidateSession).Delete("/user/{id}", userHandler.HandleDeleteUser)

	return router
}
