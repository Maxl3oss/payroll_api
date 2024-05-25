package controllers

import (
	"maxl3oss/app/models"
	"maxl3oss/pkg/response"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type DashboardController struct {
	DB *gorm.DB
}

func NewDashboardController(db *gorm.DB) *DashboardController {
	return &DashboardController{DB: db}
}

func (d *DashboardController) GetDashboard(c *fiber.Ctx) error {
	// var for search
	yearParam := c.Query("year")
	year, err := strconv.Atoi(yearParam)
	if err != nil || yearParam == "" {
		year = time.Now().Year()
	}

	// find sum
	var totalReceived float64
	var totalUser int64

	if result := d.DB.Model(&models.Salary{}).Where("EXTRACT(YEAR FROM created_at) = ?", year).Select("SUM(received) as total_received").Scan(&totalReceived); result.Error != nil {
		// return response.Message(c, fiber.ErrBadRequest.Code, false, result.Error.Error())
		totalReceived = 0
	}

	if result := d.DB.Model(&models.User{}).Where("deleted_at IS NULL AND role_id = ?", 2).Count(&totalUser); result.Error != nil {
		// return response.Message(c, fiber.StatusInternalServerError, false, result.Error.Error())
		totalUser = 0
	}

	// query by month
	// Query to sum Received field grouped by month and year
	var dataReceived = []models.ReceivedByMonth{}
	result := d.DB.Model(&models.Salary{}).
		Select("EXTRACT(YEAR FROM created_at) as year, EXTRACT(MONTH FROM created_at) as month, SUM(received) as sum").
		Where("EXTRACT(YEAR FROM created_at) = ?", year).
		Group("EXTRACT(YEAR FROM created_at), EXTRACT(MONTH FROM created_at)").
		Scan(&dataReceived)
	if result.Error != nil {
		return response.Message(c, fiber.StatusInternalServerError, false, result.Error.Error())
	}

	var data = &models.Dashboard{
		Received:        totalReceived,
		User:            totalUser,
		ReceivedByMonth: dataReceived,
	}

	return response.SendData(c, fiber.StatusOK, true, data, nil)
}
