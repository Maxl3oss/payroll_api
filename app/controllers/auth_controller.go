package controllers

import (
	"log"
	"maxl3oss/app/models"
	"maxl3oss/pkg/response"
	"maxl3oss/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type AuthController struct {
	DB *gorm.DB
}

func NewAuthController(db *gorm.DB) *AuthController {
	return &AuthController{DB: db}
}

func (a *AuthController) Register(ctx *fiber.Ctx) error {
	var user models.User

	if err := ctx.BodyParser(&user); err != nil {
		return err
	}

	// check user
	prevUser := a.DB.Where(&models.User{Email: user.Email}).First(&models.User{})
	if prevUser.Error == nil {
		return response.Message(ctx, fiber.StatusBadRequest, false, "Email already exists")
	}

	// hash password
	user.Password = utils.GeneratePassword(user.Password)
	user.RoleID = 2

	result := a.DB.Create(&user)
	if result.Error != nil {
		return response.Message(ctx, fiber.StatusInternalServerError, false, "Failed to register user")
	}

	return response.Message(ctx, fiber.StatusOK, true, "User registered successfully!")
}

func (a *AuthController) Login(ctx *fiber.Ctx) error {
	var input models.LoginInput
	if err := ctx.BodyParser(&input); err != nil {
		return err
	}

	if input.Email == "" || input.Email == "-" {
		input.Email = "XxX-xXx-XXX-xxx"
	}

	// data user from database
	var user models.User

	err := a.DB.Preload("Role").
		Where(&models.User{Email: input.Email}).
		Where("deleted_at IS NULL").
		Or(&models.User{TaxID: input.Email}).
		First(&user).Error

	if err != nil {
		return response.Message(ctx, fiber.StatusBadRequest, false, "ไม่พบผู้ใช้งานนี้!")
	}

	// Compare hashed passwords
	if !utils.ComparePasswords(user.Password, input.Password) {
		return response.Message(ctx, fiber.StatusBadRequest, false, "อีเมลหรือรหัสไม่ถูกต้อง!")
	}

	// Generate JWT token
	id := user.ID
	token, err := utils.GenerateNewTokens(id, []string{user.Role.Name})

	if err != nil {
		return response.Message(ctx, fiber.StatusInternalServerError, false, "สร้าง token ไม่สำเร็จ!")
	}

	// Return token in response
	user.Password = ""
	data := fiber.Map{"token": token, "user": user}

	return response.SendData(ctx, fiber.StatusOK, true, data, nil)
}

func (a *AuthController) RefreshToken(ctx *fiber.Ctx) error {
	var dataToken models.ReqToken
	if err := ctx.BodyParser(&dataToken); err != nil {
		return err
	}

	extractToken, err := utils.ExtractTokenMetadata(ctx)
	if err != nil {
		return response.Message(ctx, fiber.StatusInternalServerError, false, err.Error())
	}
	// Log extractToken
	userId := extractToken.UserID
	log.Printf("Extracted token metadata: %+s\n", userId)

	// refreshToken, err := utils.ParseRefreshToken(dataToken.Refresh)
	// if err != nil {
	// 	return response.Message(ctx, fiber.StatusInternalServerError, false, err.Error())
	// }
	// Log refreshToken
	// log.Printf("Parsed refresh token: %d\n", refreshToken)

	var user models.User
	result := a.DB.Preload("Prefix").First(&user, userId)
	if result.Error != nil {
		return response.Message(ctx, fiber.StatusBadRequest, false, result.Error.Error())
	}

	token, err := utils.GenerateNewTokens(userId, []string{user.Role.Name})
	if err != nil {
		return response.Message(ctx, fiber.StatusInternalServerError, false, "Failed to generate token")
	}
	// Return token in response
	data := fiber.Map{"token": token}

	return response.SendData(ctx, fiber.StatusOK, true, data, nil)
}
