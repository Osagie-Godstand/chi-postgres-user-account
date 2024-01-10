package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Osagie-Godstand/chi-postgres-user-account/internal/models"
)

func TestHandlePostUser(t *testing.T) {
	mockRepo := &MockUserRepository{}
	userHandler := NewUserHandler(mockRepo)

	userParams := models.CreateUserParams{
		Email:     "osagie@gg.uk",
		FirstName: "Osagie",
		LastName:  "Godstand",
		Password:  "Password",
	}

	jsonData, err := json.Marshal(userParams)
	if err != nil {
		t.Fatalf("Failed to marshal user params: %v", err)
	}

	req, err := http.NewRequest("POST", "/create-user", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	resRecorder := httptest.NewRecorder()

	userHandler.HandlePostUser(resRecorder, req)

	if resRecorder.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, resRecorder.Code)
	}

	var createdUser models.User
	err = json.NewDecoder(resRecorder.Body).Decode(&createdUser)
	if err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

}
