package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	_ "github.com/windeesel365/swagger-go-api/docs"
	"github.com/windeesel365/swagger-go-api/handlers"

	echoSwagger "github.com/swaggo/echo-swagger"
)

func initDB() {

	//ต้องload environment variables จาก .env file ก่อนเสมอ
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("environment variable DATABASE_URL is required")
	}

	var err error
	handlers.Db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	if err := handlers.Db.AutoMigrate(&handlers.Shopper{}); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}
}

func main() {

	initDB()

	//echo instance กับ configure middleware
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/shoppers", handlers.CreateShopperHandler)

	// GETแบบทั้งหมด หรือตาม username ก็ได้
	e.GET("/shoppers", handlers.GetAllShoppers)
	e.GET("/shoppers/:username", handlers.GetShopperByUsername)
	e.PUT("/shoppers/:username", handlers.UpdateShopperByUsername)
	e.DELETE("/shoppers/:username", handlers.DeleteShopperByUsername)

	//ใส่ swagger route
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT environment variable not set.")
	}

	//graceful shutdown //start server in goroutine
	go func() {
		if err := e.Start(":" + port); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	// รอ interrupt signal เพื่อ gracefully shutdown server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	e.Logger.Print("shutting down the server")

	// context เพื่อ timeout shutdown after 10 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// shutdown server
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
