package handlers

import (
	"bookstore/database"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type App struct {
	engine *gin.Engine

	secretKey []byte

	bh     *BookHandler
	uh     *UserHandler
	rh     *ReviewHandler
	ch     *CartHandler
}

type Config struct {
	SecretKey []byte

	BooksStorage *database.BooksStorage
	UsersStorage *database.UsersStorage
	ReviewsStorage *database.ReviewsStorage
	CartsStorage *database.CartsStorage
}

func NewApp(config *Config) *App {
	app := &App{
		secretKey: config.SecretKey,
		bh: newBookHandler(config),
		uh: newUserHandler(config),
		rh: newReviewHandler(config),
		ch: newCartHandler(config),
	}

	app.engine = registerApp(app)

	return app
}

func (app *App) pong() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "pong"})
	}
}

func (app *App) Run(hostname, port string) {
	log.Fatal(app.engine.Run(hostname + ":" + port))
}
