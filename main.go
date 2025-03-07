package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/kusnadi8605/news/config"
	"github.com/kusnadi8605/news/handler"
	"github.com/kusnadi8605/news/repository"
	"github.com/kusnadi8605/news/usecase"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func initLogger() {
	logrus.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05.999",
	})

}

func main() {
	initLogger()
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		logrus.Println("No .env file found")
	}

	// Connect to Database and Redis
	config.ConnectDB()
	config.ConnectRedis()

	// Initialize Repository, Usecase, and Handler
	newsRepo := repository.NewNewsRepository(config.DB)
	newsUsecase := usecase.NewNewsUsecase(newsRepo)
	newsHandler := handler.NewNewsHandler(newsUsecase)

	// Initialize Echo instance
	e := echo.New()

	// Define Routes
	e.GET("/news", newsHandler.GetAllNews)
	e.GET("/news/:id", newsHandler.GetNewsByID)
	e.POST("/news", newsHandler.CreateNews)
	e.PUT("/news/:id", newsHandler.UpdateNews)
	e.DELETE("/news/:id", newsHandler.DeleteNews)

	// Start Server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Starting server on port", port)
	e.Logger.Fatal(e.Start(":" + port))
}
