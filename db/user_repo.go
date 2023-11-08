package db

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/Osagie-Godstand/chi-postgres-user-account/models"
	"github.com/google/uuid"
)

type UserRepository interface {
	InsertUser(user *models.User) (*models.User, error)
	UpdateUser(userID uuid.UUID, params models.UpdateUserParams) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	GetUserByID(userID uuid.UUID) (*models.User, error)
	GetUsers() ([]models.User, error)
	DeleteUser(userID uuid.UUID) error
}

type UserPostgresRepository struct {
	DB *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) UserRepository {
	return &UserPostgresRepository{DB: db}
}

func (ur *UserPostgresRepository) InsertUser(user *models.User) (*models.User, error) {
	// Generates a new UUID for the user ID
	userID := models.NewUUID()

	// Inserting user into the database
	_, err := ur.DB.Exec(
		"INSERT INTO users (id, first_name, last_name, email, encrypted_password, is_admin) VALUES ($1, $2, $3, $4, $5, $6)",
		userID, user.FirstName, user.LastName, user.Email, user.EncryptedPassword, false,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to insert user into database: %w", err)
	}

	return user, nil
}

func (ur *UserPostgresRepository) UpdateUser(userID uuid.UUID, params models.UpdateUserParams) (*models.User, error) {
	_, err := ur.DB.Exec(
		"UPDATE users SET first_name = $1, last_name = $2 WHERE id = $3",
		params.FirstName, params.LastName, userID,
	)
	if err != nil {
		return nil, err
	}

	return ur.GetUserByID(userID)
}

func (ur *UserPostgresRepository) GetUserByEmail(email string) (*models.User, error) {
	row := ur.DB.QueryRow("SELECT id, first_name, last_name, email, encrypted_password, is_admin FROM users WHERE email = $1", email)

	var user models.User
	err := row.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.EncryptedPassword, &user.IsAdmin)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

func (ur *UserPostgresRepository) GetUserByID(userID uuid.UUID) (*models.User, error) {
	row := ur.DB.QueryRow("SELECT id, first_name, last_name, email, is_admin FROM users WHERE id = $1", userID)

	var user models.User
	if err := row.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.IsAdmin); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user with ID %s not found", userID.String())
		}
		return nil, err
	}

	return &user, nil
}

func (ur *UserPostgresRepository) GetUsers() ([]models.User, error) {
	rows, err := ur.DB.Query("SELECT id, first_name, last_name, email, is_admin FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.IsAdmin); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (ur *UserPostgresRepository) DeleteUser(userID uuid.UUID) error {
	_, err := ur.DB.Exec("DELETE FROM users WHERE id = $1", userID)
	return err
}
