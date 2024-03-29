package controllers

import (
	"maxl3oss/app/models"
	"maxl3oss/pkg/response"
	"maxl3oss/pkg/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type UserController struct {
	DB *gorm.DB
}

func NewUserController(db *gorm.DB) *UserController {
	return &UserController{DB: db}
}

func (u *UserController) GetAll(ctx *fiber.Ctx) error {
	var users []models.User

	result := u.DB.Model(&models.User{}).Preload("Role").Find(&users)
	if result.Error != nil {
		return response.Message(ctx, fiber.StatusInternalServerError, false, result.Error.Error())
	}

	return response.SendData(ctx, fiber.StatusOK, true, users, nil)
}

func (u *UserController) Create(ctx *fiber.Ctx) error {
	var user models.User

	if err := ctx.BodyParser(&user); err != nil {
		return err
	}

	// hash password
	user.Password = utils.GeneratePassword(user.Password)

	result := u.DB.Create(&user)
	if result.Error != nil {
		return result.Error
	}

	return ctx.JSON("Create successfully!")
}

func (u *UserController) DeleteById(ctx *fiber.Ctx) error {
	idStr := ctx.Params("id")

	// Convert idStr to integer
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return err
	}

	// find
	check := u.DB.First(&models.User{}, id)
	if check.Error != nil {
		return check.Error
	}

	result := u.DB.Delete(&models.User{}, id)
	if result.Error != nil {
		return result.Error
	}

	return ctx.JSON("Delete successfully!")
}

func (u *UserController) Update(ctx *fiber.Ctx) error {

	// Get ID from parameters
	id := ctx.Params("id")

	// Find user by ID
	var user models.User
	if result := u.DB.First(&user, id); result.Error != nil {
		return response.Message(ctx, fiber.StatusNotFound, false, result.Error.Error())
	}

	// Parse updated fields from request body
	type UpdateUser struct {
		FullName string
		Email    string
	}
	var updateData UpdateUser
	if err := ctx.BodyParser(&updateData); err != nil {
		return response.Message(ctx, fiber.StatusBadRequest, false, err.Error())
	}

	// Update fields
	user.FullName = updateData.FullName
	user.Email = updateData.Email

	// Save changes
	result := u.DB.Save(&user)
	if result.Error != nil {
		return response.Message(ctx, fiber.StatusInternalServerError, false, result.Error.Error())
	}

	// Return response
	return response.Message(ctx, fiber.StatusOK, true, "User updated successfully")
}

func (u *UserController) GetAllPrefix(ctx *fiber.Ctx) error {
	var prefix []models.Prefix

	result := u.DB.Find(&prefix)
	if result.Error != nil {
		return result.Error
	}

	return ctx.JSON(prefix)
}
