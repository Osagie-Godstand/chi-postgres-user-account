package handlers

import (
	"github.com/Osagie-Godstand/chi-postgres-user-account/internal/models"
	"github.com/google/uuid"
)

type MockUserRepository struct{}

func (m *MockUserRepository) InsertUser(user *models.User) (*models.User, error) {
	return &models.User{
		ID:        uuid.New(),
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}, nil
}

func (m *MockUserRepository) DeleteUser(userID uuid.UUID) error {
	return nil
}

func (m *MockUserRepository) GetUserByEmail(email string) (*models.User, error) {
	return &models.User{
		Email: email,
	}, nil
}

func (m *MockUserRepository) GetUserByID(userID uuid.UUID) (*models.User, error) {
	return &models.User{
		ID: userID,
	}, nil
}

func (m *MockUserRepository) GetUsers() ([]models.User, error) {
	users := []models.User{
		{
			ID:        uuid.New(),
			Email:     "user1@example.com",
			FirstName: "User",
			LastName:  "One",
		},
		{
			ID:        uuid.New(),
			Email:     "user2@example.com",
			FirstName: "User",
			LastName:  "Two",
		},
	}
	return users, nil
}

func (m *MockUserRepository) UpdateUser(userID uuid.UUID, params models.UpdateUserParams) (*models.User, error) {
	updatedUser := &models.User{
		ID:        userID,
		FirstName: params.FirstName,
		LastName:  params.LastName,
	}
	return updatedUser, nil
}
