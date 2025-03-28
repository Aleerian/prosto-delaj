package repository

import (
	"database/sql"
	"github.com/sirupsen/logrus"
	"prosto-delaj-api/models"
	"prosto-delaj-api/server/queries"
)

func NewBusinessDatabase(config *models.BusinessDBConfig) (*sql.DB, *queries.Queries) {
	logrus.Info("start database connected")
	database, err := NewPostgresDB(&models.BusinessDBConfig{
		Host:     config.Host,
		Port:     config.Port,
		Username: config.Username,
		Password: config.Password,
		DBName:   config.DBName,
		SSLMode:  config.SSLMode,
	})
	if err != nil {
		logrus.Fatalf("failed to initialize business db: %s", err.Error())
	}
	logrus.Info("database connected")
	db := queries.New(database)
	logrus.Info("queries connected")
	return database, db
}
