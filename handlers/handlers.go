package handlers

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	_ "github.com/windeesel365/swagger-go-api/docs"
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

var Db *gorm.DB

// createShopperHandler godoc
// @Summary Create a new shopper
// @Description Create a new shopper with the provided data
// @Accept json
// @Produce json
// @Param shopper body Shopper true "Shopper object to be created"
// @Success 201 {object} Shopper "Created shopper object"
// @Failure 400 {object} ErrorResponse "Error response"
// @Router /shoppers [post]
func CreateShopperHandler(c echo.Context) error {
	shopper := new(Shopper)
	if err := c.Bind(shopper); err != nil {
		return err
	}

	shopper.DateJoined = time.Now().Format("2006-01-02")

	//save to database
	if err := Db.Create(&shopper).Error; err != nil {
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
func GetAllShoppers(c echo.Context) error {
	var shoppers []Shopper
	if err := Db.Find(&shoppers).Error; err != nil {
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
func GetShopperByUsername(c echo.Context) error {
	username := c.Param("username")
	var shopper Shopper
	if err := Db.Where("username = ?", username).First(&shopper).Error; err != nil {
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
func UpdateShopperByUsername(c echo.Context) error {
	username := c.Param("username")
	var updatedShopper Shopper
	if err := c.Bind(&updatedShopper); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request payload"})
	}
	var existingShopper Shopper
	if err := Db.Where("username = ?", username).First(&existingShopper).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "shopper not found"})
	}
	existingShopper.FullName = updatedShopper.FullName
	existingShopper.Email = updatedShopper.Email
	existingShopper.Street = updatedShopper.Street
	existingShopper.City = updatedShopper.City
	existingShopper.State = updatedShopper.State
	existingShopper.ZipCode = updatedShopper.ZipCode

	if err := Db.Save(&existingShopper).Error; err != nil {
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
func DeleteShopperByUsername(c echo.Context) error {
	username := c.Param("username")
	var shopper Shopper
	if err := Db.Where("username = ?", username).First(&shopper).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "shopper not found"})
	}
	if err := Db.Delete(&shopper).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to delete shopper"})
	}
	return c.NoContent(http.StatusNoContent)
}
