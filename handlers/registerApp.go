package handlers

import (
	"bookstore/middleware"

	"github.com/gin-gonic/gin"
)

func registerApp(app *App) *gin.Engine {
	engine := gin.Default()

	engine.GET("/ping", app.pong())

	registerBookHandlers(engine, app.bh, app.secretKey)

	registerUserHandlers(engine, app.uh, app.secretKey)

	registerReviewHandlers(engine, app.rh, app.secretKey)

	registerCartHandlers(engine, app.ch, app.secretKey)

	registerNewsHandlers(engine, app.nh)

	registerStatic(engine)

	return engine
}

func registerBookHandlers(engine *gin.Engine, bh *BookHandler, secretKey []byte) {
	engine.GET("/book/:id", bh.GetBook())

	engine.GET("/book/all", bh.GetAllBooks())

	engine.GET("/book", bh.GetTop())

	engine.GET("/book/search", bh.SearchBooks())

	group := engine.Group("/book/admin", middleware.ValidateToken(secretKey))
	{
		group.POST("/", bh.AddBook())
		group.DELETE("/:id", bh.DeleteBook())
	}
}

func registerUserHandlers(engine *gin.Engine, handler *UserHandler, secretKey []byte) {
	group := engine.Group("/user", middleware.ValidateToken(secretKey))
	{
		group.GET("/:id", handler.GetUser())
		group.DELETE("/:id", handler.DeleteUser())
		group.POST("/", handler.AddUser())

		group.GET("/all", handler.GetAllUsers())

		group.GET("/search", handler.GetUserByUsername())

		group.GET("/refresh", handler.RefreshToken())
	}

	engine.POST("/login", handler.Login())
	engine.POST("/signup", handler.SignUp())
}

func registerReviewHandlers(engine *gin.Engine, handler *ReviewHandler, secretKey []byte) {
	group := engine.Group("/review", middleware.ValidateToken(secretKey))
	{
		group.DELETE("/:id", handler.DeleteReview())
		group.POST("/", handler.AddReview())
	}

	engine.GET("/review/:id", handler.GetReview())

	engine.GET("/reviews/all", handler.GetAllReviews())

	engine.GET("/review/searchby/user", handler.GetReviewsByUserId())
	engine.GET("/review/searchby/book", handler.GetReviewsByBookId())
}

func registerCartHandlers(engine *gin.Engine, handler *CartHandler, secretKey []byte) {
	group := engine.Group("/cart", middleware.ValidateToken(secretKey))
	{
		group.GET("/:id", handler.GetCart())
		group.DELETE("/:id", handler.DeleteCart())
		group.POST("/", handler.AddCart())

		group.GET("/all", handler.GetAllCarts())

		group.GET("/searchby/user", handler.GetCartsByUserId())
		group.GET("/searchby/book", handler.GetCartsByBookId())
	}
}

func registerNewsHandlers(engine *gin.Engine, handler *NewsHandler) {
	group := engine.Group("/news")
	{
		group.GET("/:id", handler.GetNews())

		group.GET("/all", handler.GetAllNews())
	}
}

func registerStatic(engine *gin.Engine) {
	engine.StaticFile("/", "./public/index.html")
	engine.Static("/scripts", "./public")
	engine.Static("/pdf", "./public/storage/pdf")
	engine.Static("/cover", "./public/storage/cover")
	engine.Static("/video", "./public/storage/video")
	engine.Static("/logo", "./public/storage/logo")
	engine.Static("/newspic", "./public/storage/newspic")
}
