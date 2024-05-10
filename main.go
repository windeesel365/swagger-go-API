package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	echoSwagger "github.com/swaggo/echo-swagger"
)

type Shopper struct {
	Username   string `json:"username" gorm:"primaryKey"`
	FullName   string `json:"fullName"`
	Email      string `json:"email"`
	Street     string `json:"street"`
	City       string `json:"city"`
	State      string `json:"state"`
	ZipCode    string `json:"zipCode"`
	DateJoined string `json:"dateJoined"`
}

type ShoppersResponse struct {
	Shoppers []Shopper `json:"shoppers"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

var db *gorm.DB

func initDB() {

	//ต้อง load environment variables จาก .env file ก่อนเสมอ
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("environment variable DATABASE_URL is required")
	}

	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	if err := db.AutoMigrate(&Shopper{}); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}
}

func createShopperHandler(c echo.Context) error {
	shopper := new(Shopper)
	if err := c.Bind(shopper); err != nil {
		return err
	}

	shopper.DateJoined = time.Now().Format("2006-01-02")

	// Save to database
	if err := db.Create(&shopper).Error; err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, shopper)
}

func main() {

	initDB()

	//echo instance กับ configure middleware
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/shoppers", createShopperHandler)

	// GETแบบทั้งหมด หรือตาม username ก็ได้
	e.GET("/shoppers", getAllShoppers)
	e.GET("/shoppers/:username", getShopperByUsername)
	e.PUT("/shoppers/:username", updateShopperByUsername)
	e.DELETE("/shoppers/:username", deleteShopperByUsername)

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
