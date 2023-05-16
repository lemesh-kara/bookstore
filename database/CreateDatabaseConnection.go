package database

import (
	"strings"

	"bookstore/logger"

	"github.com/jinzhu/gorm"
)

func CreateDatabaseConnection(uri, username, password, dbname string) *gorm.DB {
	connectUri := "postgres://" + username + ":" + password + "@" + uri + "/" + dbname + "?sslmode=disable"

	db, err := gorm.Open("postgres", connectUri)
	if err != nil {
		logger.Log.Fatal("Cannot connect to database: " + err.Error())
	}

	logger.Log.Info("Connected to database " + uri)

	return db
}

func IsNotFoundError(errorString string) bool {
	return strings.Contains(errorString, "record not found")
}
