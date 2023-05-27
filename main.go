package main

import (
	"bookstore/database"
	"bookstore/handlers"
	"bookstore/logger"
	"bookstore/testdata"
	"os"
)

func main() {
	logger.InitLogger()

	logger.Log.Info("Starting application")

	db := database.CreateDatabaseConnection("localhost:5432", "username", "password", "bookstore")

	testdata.InitDatabase(db)
	testdata.CreateTestData(db)

	var config handlers.Config
	config.BooksStorage = database.NewBooksStorage(db)
	config.UsersStorage = database.NewUsersStorage(db)
	config.ReviewsStorage = database.NewReviewsStorage(db)
	config.CartsStorage = database.NewCartsStorage(db)
	config.NewsStorage = database.NewNewsStorage(db)

	config.SecretKey = []byte(os.Getenv("SECRET_KEY"))

	logger.Log.Debug("Using secret key " + string(config.SecretKey))

	app := handlers.NewApp(&config)

	app.Run("", "8080")
}
