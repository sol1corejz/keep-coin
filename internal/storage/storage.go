package storage

import (
	"database/sql"
	"errors"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/sol1corejz/keep-coin/cmd/config"
	"github.com/sol1corejz/keep-coin/internal/logger"
	"github.com/sol1corejz/keep-coin/internal/models"
	"go.uber.org/zap"
)

var DB *sql.DB
var ErrConnectionFailed = errors.New("db connection failed")

func Init() error {
	if config.DatabaseURI == "" {
		logger.Log.Error("No DB URI")
		return ErrConnectionFailed
	}

	db, err := sql.Open("pgx", config.DatabaseURI)
	if err != nil {
		logger.Log.Fatal("Error opening database connection", zap.Error(err))
		return ErrConnectionFailed
	}
	DB = db

	tables := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY NOT NULL,
			first_name VARCHAR(255) DEFAULT '',
			last_name VARCHAR(255) DEFAULT '',
			email VARCHAR(255) UNIQUE NOT NULL,
			password VARCHAR(255) NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);`,
	}

	for _, table := range tables {
		if _, err := DB.Exec(table); err != nil {
			errors.New("creating table failed")
		}
	}

	return nil
}

func RegisterUser(user *models.User) error {

	_, err := DB.Exec(`
		INSERT INTO users (id, email, password) 
		VALUES ($1, $2, $3);`, user.Uuid, user.Email, user.Password)

	if err != nil {
		logger.Log.Error("Error registering user", zap.Error(err))
		return errors.New("creating user failed")
	}

	return nil
}

func GetUserByLogin(email string) (models.User, error) {

	var existingUser models.User

	err := DB.QueryRow(`
		SELECT * FROM users WHERE email = $1;
	`, email).Scan(&existingUser.Uuid, &existingUser.FirstName, &existingUser.LastName, &existingUser.Email, &existingUser.Password, &existingUser.CreatedAt)

	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return models.User{}, err
		}
		return models.User{}, err
	}

	return existingUser, nil
}
