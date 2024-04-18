package controllers

import (
	"maxl3oss/app/models"
	"maxl3oss/pkg/response"
	"maxl3oss/pkg/utils"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type UserController struct {
	DB *gorm.DB
}

func NewUserController(db *gorm.DB) *UserController {
	return &UserController{DB: db}
}

func (u *UserController) GetAll(c *fiber.Ctx) error {
	// Define variables to hold users and result count
	var users []models.User
	var resCount int64

	// Get query parameters for pagination and search
	page, _ := strconv.Atoi(c.Query("pageNumber", "1"))
	limit, _ := strconv.Atoi(c.Query("pageSize", "10"))
	search := c.Query("search")

	// Calculate offset for pagination
	offset := (page - 1) * limit

	// Query to get the total count of users based on search criteria
	queryCount := u.DB.Model(&models.User{}).Where("deleted_at IS NULL AND full_name LIKE ?", "%"+search+"%").Or("email LIKE ?", "%"+search+"%")
	err := queryCount.Count(&resCount).Error
	if err != nil {
		return response.Message(c, fiber.StatusInternalServerError, false, err.Error())
	}

	// Query to retrieve users based on role and search criteria, with pagination
	result := u.DB.Select("id", "full_name", "email", "mobile", "role_id", "created_at").
		Preload("Role").
		Where("deleted_at IS NULL AND full_name LIKE ?", "%"+search+"%").Or("email LIKE ?", "%"+search+"%").
		Order("created_at ASC").
		Limit(limit).Offset(offset).
		Find(&users)

	// Check for errors in the database query
	if result.Error != nil {
		return response.Message(c, fiber.StatusInternalServerError, false, result.Error.Error())
	}

	// Create pagination object
	pagin := response.Pagination{
		PageNumber:  page,
		PageSize:    limit,
		TotalRecord: int(resCount),
	}

	// Send response with user data and pagination information
	return response.SendData(c, fiber.StatusOK, true, users, &pagin)
}

func (u *UserController) GetByUser(c *fiber.Ctx) error {
	userID := c.Params("id")
	var user models.User
	// check user
	check := u.DB.Where("id = ?", userID).First(&user)
	if check.Error != nil {
		return response.Message(c, fiber.StatusInternalServerError, false, check.Error.Error())
	}

	return response.SendData(c, fiber.StatusOK, true, user, nil)
}

func (u *UserController) GetProfile(c *fiber.Ctx) error {
	// userID := c.Params("id")
	extractToken, err := utils.ExtractTokenMetadata(c)
	if err != nil {
		// Handle error
		return response.Message(c, 401, false, "ไม่สามารถเข้าถึงข้อมูล Token ได้")
	}
	var user models.User

	// check user
	check := u.DB.Preload("Role").Where("id = ?", extractToken.UserID).First(&user)
	if check.Error != nil {
		return response.Message(c, fiber.StatusInternalServerError, false, check.Error.Error())
	}

	return response.SendData(c, fiber.StatusOK, true, user, nil)
}

func (u *UserController) Create(c *fiber.Ctx) error {
	var user models.User

	if err := c.BodyParser(&user); err != nil {
		return err
	}

	// hash password
	user.Password = utils.GeneratePassword(user.Password)

	result := u.DB.Create(&user)
	if result.Error != nil {
		return response.Message(c, fiber.StatusNotFound, false, result.Error.Error())
	}

	return response.Message(c, fiber.StatusOK, true, "Create successfully!")
}

func (u *UserController) DeleteById(c *fiber.Ctx) error {
	userID := c.Params("id")
	// Find user by ID
	var user models.User
	if result := u.DB.Where("id = ?", userID).First(&user); result.Error != nil {
		return response.Message(c, fiber.StatusNotFound, false, result.Error.Error())
	}

	// update
	now := time.Now()
	user.DeletedAt = &now

	// Save changes
	if result := u.DB.Save(&user); result.Error != nil {
		return response.Message(c, fiber.StatusInternalServerError, false, result.Error.Error())
	}
	// result := u.DB.Model(&models.User{}).Delete("id = ?", userID)
	// if result.Error != nil {
	// 	return response.Message(c, fiber.StatusBadRequest, false, result.Error.Error())
	// }

	return response.Message(c, fiber.StatusOK, true, "Delete successfully!")
}

func (u *UserController) Update(ctx *fiber.Ctx) error {

	// Get ID from parameters
	uid := ctx.Params("id")

	// Find user by ID
	var user models.User
	if result := u.DB.Model(&models.User{}).Where("id = ?", uid).First(&user); result.Error != nil {
		return response.Message(ctx, fiber.StatusNotFound, false, result.Error.Error())
	}

	// Parse updated fields from request body
	type UpdateUser struct {
		FullName string `json:"full_name"`
		Email    string `json:"email"`
		// PassWord string `json:"password"`
		RoleId int16  `json:"role_id"`
		Mobile string `json:"mobile"`
	}

	var updateData UpdateUser
	if err := ctx.BodyParser(&updateData); err != nil {
		return response.Message(ctx, fiber.StatusBadRequest, false, err.Error())
	}

	// Update fields
	user.FullName = updateData.FullName
	user.Email = updateData.Email
	user.Mobile = updateData.Mobile
	user.RoleID = updateData.RoleId
	// if user.Password != "" {
	// 	user.Password = utils.GeneratePassword(updateData.PassWord)
	// }

	// Save changes
	result := u.DB.Save(&user)
	if result.Error != nil {
		return response.Message(ctx, fiber.StatusInternalServerError, false, result.Error.Error())
	}

	// Return response
	return response.Message(ctx, fiber.StatusOK, true, "User updated successfully")
}

func (u *UserController) ChangePassByAdmin(ctx *fiber.Ctx) error {

	// Get ID from parameters
	uid := ctx.Params("id")

	// Find user by ID
	var user models.User
	if result := u.DB.Model(&models.User{}).Where("id = ?", uid).First(&user); result.Error != nil {
		return response.Message(ctx, fiber.StatusNotFound, false, result.Error.Error())
	}

	// Parse updated fields from request body
	type UpdatePass struct {
		ConfirmPassword string `json:"confirm_password"`
		PassWord        string `json:"password"`
	}

	var updateData UpdatePass
	if err := ctx.BodyParser(&updateData); err != nil {
		return response.Message(ctx, fiber.StatusBadRequest, false, err.Error())
	}

	// check pass
	if updateData.PassWord != updateData.ConfirmPassword {
		return response.Message(ctx, fiber.StatusInternalServerError, false, "รหัสผ่านไม่ตรงกัน!")
	}

	// Update fields
	if user.Password != "" {
		user.Password = utils.GeneratePassword(updateData.PassWord)
	}

	// Save changes
	result := u.DB.Save(&user)
	if result.Error != nil {
		return response.Message(ctx, fiber.StatusInternalServerError, false, result.Error.Error())
	}

	// Return response
	return response.Message(ctx, fiber.StatusOK, true, "User updated successfully")
}

func (u *UserController) ChangePassByUser(ctx *fiber.Ctx) error {

	// Get ID from parameters
	uid := ctx.Params("id")

	// Find user by ID
	var user models.User
	if result := u.DB.Model(&models.User{}).Where("id = ?", uid).First(&user); result.Error != nil {
		return response.Message(ctx, fiber.StatusNotFound, false, result.Error.Error())
	}

	// Parse updated fields from request body
	type UpdatePass struct {
		PrevPassword    string `json:"prev_password"`
		ConfirmPassword string `json:"confirm_password"`
		PassWord        string `json:"password"`
	}

	var updateData UpdatePass
	if err := ctx.BodyParser(&updateData); err != nil {
		return response.Message(ctx, fiber.StatusBadRequest, false, err.Error())
	}

	// Compare hashed passwords
	if !utils.ComparePasswords(user.Password, updateData.PrevPassword) {
		return response.Message(ctx, fiber.StatusBadRequest, false, "รหัสผ่านไม่ถูกต้อง!")
	}

	// check pass
	if updateData.PassWord != updateData.ConfirmPassword {
		return response.Message(ctx, fiber.StatusInternalServerError, false, "รหัสผ่านไม่ตรงกัน!")
	}

	// Update fields
	if user.Password != "" {
		user.Password = utils.GeneratePassword(updateData.PassWord)
	}

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

// role
func (u *UserController) GetRole(c *fiber.Ctx) error {
	var role []models.Role

	if result := u.DB.Model(&models.Role{}).Find(&role); result.Error != nil {
		return response.Message(c, fiber.StatusBadRequest, false, result.Error.Error())
	}

	return response.SendData(c, fiber.StatusOK, true, role, nil)
}
