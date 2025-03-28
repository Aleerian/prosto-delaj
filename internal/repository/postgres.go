package repository

import (
	"database/sql"
	"fmt"
	"prosto-delaj-api/models"

	_ "github.com/jackc/pgx/stdlib"
)

const (
	dbDriverName = "pgx"
)

func NewPostgresDB(config *models.BusinessDBConfig) (*sql.DB, error) {

	db, err := getDBConnection(config)
	if err == nil {
		return db, nil
	}

	return nil, err
}

func getDBConnection(config *models.BusinessDBConfig) (*sql.DB, error) {
	return sql.Open(
		dbDriverName, fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			config.Host, config.Port, config.Username, config.Password, config.DBName, config.SSLMode))
}
