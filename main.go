package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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

}
