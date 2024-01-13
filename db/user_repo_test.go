package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/DATA-DOG/go-txdb"
	"github.com/Osagie-Godstand/chi-postgres-user-account/internal/models"
	_ "github.com/lib/pq"
	migrate "github.com/rubenv/sql-migrate"
)

func TestUserRepository(t *testing.T) {
	// Register the go-txdb driver
	txdb.Register("txdb", "postgres", "user=postgres password=postgres dbname=testdb sslmode=disable")

	// Connect to the temporary test database
	sqlDB, err := sql.Open("txdb", "testdb")
	if err != nil {
		t.Fatal(err)
	}
	defer sqlDB.Close()

	err = runMigrations(sqlDB)
	if err != nil {
		t.Fatal(err)
	}

	userPostgresRepository := &UserPostgresRepository{
		DB: sqlDB,
	}

	user := models.User{
		FirstName:         "Osagie",
		LastName:          "Godstand",
		Email:             "osagie@gg.uk",
		EncryptedPassword: "",
		IsAdmin:           false,
	}

	_, err = userPostgresRepository.InsertUser(&user)
	if err != nil {
		t.Fatalf("Error inserting user: %v", err)
	}

	retrievedUser, err := userPostgresRepository.GetUserByEmail("osagie@gg.uk")
	if err != nil {
		t.Fatalf("Error retrieving user: %v", err)
	}
	if retrievedUser.FirstName != "Osagie" || retrievedUser.LastName != "Godstand" {
		t.Fatalf("Unexpected user properties")
	}

	// other test cases can be added here

}

func runMigrations(db *sql.DB) error {
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting current directory: %v", err)
	}

	migrationsDir := filepath.Join(currentDir, "migrations_test")

	migrations := &migrate.FileMigrationSource{
		Dir: migrationsDir,
	}

	if err := dropTables(db); err != nil {
		return fmt.Errorf("error dropping existing tables: %v", err)
	}

	n, err := migrate.Exec(db, "postgres", migrations, migrate.Up)
	if err != nil {
		return fmt.Errorf("error applying migrations: %v", err)
	}

	log.Printf("Applied %d migrations!\n", n)

	return nil
}

func dropTables(db *sql.DB) error {
	statements := []string{
		"DROP TABLE IF EXISTS users;",
	}

	for _, statement := range statements {
		_, err := db.Exec(statement)
		if err != nil {
			return err
		}
	}

	return nil
}
