package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Osagie-Godstand/chi-postgres-user-account/db"
	"github.com/Osagie-Godstand/chi-postgres-user-account/models"
	"github.com/google/uuid"
)

type SessionHandler struct {
	DB             *sql.DB
	userRepository db.UserRepository
}

func NewSessionHandler(db *sql.DB, userRepository db.UserRepository) *SessionHandler {
	return &SessionHandler{
		DB:             db,
		userRepository: userRepository,
	}
}

type LoginParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResp struct {
	Type string `json:"type"`
	Msg  string `json:"msg"`
}

func invalidCredentials(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	resp := loginResp{
		Type: "error",
		Msg:  "invalid credentials",
	}
	err := json.NewEncoder(w).Encode(resp)
	return err
}

func (s *SessionHandler) CreateSession(session models.Session) error {
	_, err := s.DB.Exec(
		"INSERT INTO sessions (user_id, token, expires_at) VALUES ($1, $2, $3)",
		session.UserID, session.Token, session.ExpiresAt,
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *SessionHandler) GetSessionByToken(token string) (*models.Session, error) {
	row := s.DB.QueryRow("SELECT id, user_id, token, expires_at FROM sessions WHERE token = $1", token)
	var session models.Session
	if err := row.Scan(&session.ID, &session.UserID, &session.Token, &session.ExpiresAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("session not found")
		}
		return nil, err
	}
	return &session, nil
}

func (s *SessionHandler) DeleteSession(token string) error {
	_, err := s.DB.Exec("DELETE FROM sessions WHERE token = $1", token)
	return err
}

func (s *SessionHandler) GenerateSession(userID uuid.UUID) (*models.Session, error) {
	token, err := GenerateRandomToken(64)
	if err != nil {
		return nil, err
	}
	expiresAt := time.Now().Add(time.Hour * 1) // Set session expiration time 1 hour
	session := models.Session{
		UserID:    userID,
		Token:     token,
		ExpiresAt: expiresAt,
	}
	err = s.CreateSession(session)
	if err != nil {
		return nil, err
	}
	return &session, nil
}

// Login function to authenticate the user and generate a session
func (s *SessionHandler) Login(w http.ResponseWriter, r *http.Request) {
	var params LoginParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	user, err := s.userRepository.GetUserByEmail(params.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("User not found for email: %s", params.Email)
			invalidCredentials(w)
			return
		}
		log.Printf("Error fetching user by email: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !models.IsValidPassword(user.EncryptedPassword, params.Password) {
		invalidCredentials(w)
		return
	}

	session, err := s.GenerateSession(user.ID)
	if err != nil {
		log.Printf("Error generating session: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Sets the session token as a response header or in a cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    session.Token,
		Expires:  session.ExpiresAt,
		HttpOnly: true, // HttpOnly cookie important for security
		SameSite: http.SameSiteStrictMode,
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(user.ID); err != nil {
		log.Printf("Error encoding user ID as JSON response: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Logout function to delete the user's session
func (s *SessionHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// Gets the session token from the request
	token, err := r.Cookie("session_token")
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Deletes the session from the repository
	err = s.DeleteSession(token.Value)
	if err != nil {
		http.Error(w, "Failed to log out", http.StatusInternalServerError)
		return
	}

	// Expire the session cookie by setting the expiration time in the past
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Logged out successfully"})
}

// ValidateSession middleware to check if the user's session is valid
func (s *SessionHandler) ValidateSession(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the session token from the request
		token, err := r.Cookie("session_token")
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Retrieves the session from the repository
		session, err := s.GetSessionByToken(token.Value)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Checks if the session has expired
		if time.Now().After(session.ExpiresAt) {
			// Expire the session by deleting it from the repository
			err := s.DeleteSession(token.Value)
			if err != nil {
				http.Error(w, "Failed to log out", http.StatusInternalServerError)
				return
			}

			// Expire the session cookie by setting the expiration time in the past
			http.SetCookie(w, &http.Cookie{
				Name:     "session_token",
				Value:    "",
				Expires:  time.Now().Add(-time.Hour),
				HttpOnly: true,
				SameSite: http.SameSiteStrictMode,
			})

			http.Error(w, "Session expired", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "userID", session.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
