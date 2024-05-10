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

	_ "swagger-go-api/docs"

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

	//ต้องload environment variables จาก .env file ก่อนเสมอ
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

// createShopperHandler godoc
// @Summary Create a new shopper
// @Description Create a new shopper with the provided data
// @Accept json
// @Produce json
// @Param shopper body Shopper true "Shopper object to be created"
// @Success 201 {object} Shopper "Created shopper object"
// @Failure 400 {object} ErrorResponse "Error response"
// @Router /shoppers [post]
func createShopperHandler(c echo.Context) error {
	shopper := new(Shopper)
	if err := c.Bind(shopper); err != nil {
		return err
	}

	shopper.DateJoined = time.Now().Format("2006-01-02")

	//save to database
	if err := db.Create(&shopper).Error; err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, shopper)
}

// getAllShoppers godoc
// @Summary Get all shoppers
// @Description Retrieve a list of all shoppers
// @Tags shoppers
// @Produce json
// @Success 200 {object} ShoppersResponse
// @Failure 500 {object} map[string]string
// @Router /shoppers [get]
func getAllShoppers(c echo.Context) error {
	var shoppers []Shopper
	if err := db.Find(&shoppers).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "could not fetch shoppers"})
	}
	return c.JSON(http.StatusOK, ShoppersResponse{Shoppers: shoppers})
}

// getShopperByUsername godoc
// @Summary Get shopper by username
// @Description Retrieve a shopper by their username
// @Tags shoppers
// @Produce json
// @Param username path string true "Shopper Username"
// @Success 200 {object} Shopper
// @Failure 404 {object} map[string]string
// @Router /shoppers/{username} [get]
func getShopperByUsername(c echo.Context) error {
	username := c.Param("username")
	var shopper Shopper
	if err := db.Where("username = ?", username).First(&shopper).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "shopper not found"})
	}
	return c.JSON(http.StatusOK, shopper)
}

// updateShopperByUsername godoc
// @Summary Update shopper by username
// @Description Update a shopper's information by their username
// @Tags shoppers
// @Accept json
// @Produce json
// @Param username path string true "Shopper Username"
// @Param shopper body Shopper true "Shopper object to update"
// @Success 200 {object} Shopper
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /shoppers/{username} [put]
func updateShopperByUsername(c echo.Context) error {
	username := c.Param("username")
	var updatedShopper Shopper
	if err := c.Bind(&updatedShopper); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request payload"})
	}
	var existingShopper Shopper
	if err := db.Where("username = ?", username).First(&existingShopper).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "shopper not found"})
	}
	existingShopper.FullName = updatedShopper.FullName
	existingShopper.Email = updatedShopper.Email
	existingShopper.Street = updatedShopper.Street
	existingShopper.City = updatedShopper.City
	existingShopper.State = updatedShopper.State
	existingShopper.ZipCode = updatedShopper.ZipCode

	if err := db.Save(&existingShopper).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to update shopper"})
	}
	return c.JSON(http.StatusOK, existingShopper)
}

// deleteShopperByUsername godoc
// @Summary Delete shopper by username
// @Description Delete a shopper by their username
// @Tags shoppers
// @Param username path string true "Shopper Username"
// @Success 204
// @Failure 404 {object} map[string]string
// @Router /shoppers/{username} [delete]
func deleteShopperByUsername(c echo.Context) error {
	username := c.Param("username")
	var shopper Shopper
	if err := db.Where("username = ?", username).First(&shopper).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "shopper not found"})
	}
	if err := db.Delete(&shopper).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to delete shopper"})
	}
	return c.NoContent(http.StatusNoContent)
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
