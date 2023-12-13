package main

import (
	"database/sql"
	"net/http"

	"github.com/Osagie-Godstand/chi-postgres-user-account/db"
	"github.com/Osagie-Godstand/chi-postgres-user-account/handlers"
	"github.com/go-chi/chi/v5"
)

func initializeRouter(dbConn *sql.DB) http.Handler {
	router := chi.NewRouter()

	// Initializing the repositories and handlers
	userRepository := db.NewPostgresUserRepository(dbConn)
	session := handlers.NewSessionHandler(dbConn, userRepository)
	userHandler := handlers.NewUserHandler(userRepository)

	// Defining routes and handlers

	// create user and get users does not need session validation
	router.Post("/user", userHandler.HandlePostUser)

	router.Get("/users", userHandler.HandleGetUsers)

	// login is used to generate session
	router.Post("/login", session.Login)

	// Applying the ValidateSession middleware to routes that need session validation
	router.With(session.ValidateSession).Post("/logout", session.Logout)
	router.With(session.ValidateSession).Put("/user/{id}", userHandler.HandlePutUser)
	router.With(session.ValidateSession).Get("/user/{id}", userHandler.HandleGetUserByID)
	router.With(session.ValidateSession).Delete("/user/{id}", userHandler.HandleDeleteUser)

	return router
}
