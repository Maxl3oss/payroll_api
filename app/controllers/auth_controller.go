package controllers

import (
	"log"
	"maxl3oss/app/models"
	"maxl3oss/pkg/response"
	"maxl3oss/pkg/utils"
	"time"

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
		Or(&models.User{TaxID: input.Email}).
		First(&user).Error

	if user.DeletedAt != nil {
		return response.Message(ctx, fiber.StatusBadRequest, false, "บัญชีของท่านถูกปิดใช้งาน!")
	}

	if err != nil {
		return response.Message(ctx, fiber.StatusBadRequest, false, "อีเมลหรือรหัสไม่ถูกต้อง!")
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

func (a *AuthController) ForgetPassword(c *fiber.Ctx) error {
	var input struct {
		Email string `json:"email"`
	}
	if err := c.BodyParser(&input); err != nil {
		return err
	}

	// Check if email exists
	var user models.User
	if result := a.DB.Where("email = ?", input.Email).First(&user); result.Error != nil {
		return response.Message(c, fiber.StatusBadRequest, false, "ไม่พบผู้ใช้งานนี้!")
	}

	// Generate OTP
	otp := utils.GenerateOTP()
	user.OTP = otp
	user.OTPExpiry = time.Now().Add(time.Minute * 5) // OTP expiry in 5 minutes
	a.DB.Save(&user)

	// Send OTP to email
	if err := utils.SendEmail(input.Email, otp); err != nil {
		return response.Message(c, fiber.StatusBadRequest, false, err.Error())
	}

	return response.Message(c, fiber.StatusOK, true, "OTP ส่งไปในอีเมลของคุณแล้ว")
}

func (a *AuthController) ConfirmOTP(c *fiber.Ctx) error {
	var input struct {
		Email string `json:"email"`
		OTP   string `json:"otp"`
	}
	if err := c.BodyParser(&input); err != nil {
		return err
	}

	// Check if email exists
	var user models.User
	if result := a.DB.Where("email = ?", input.Email).First(&user); result.Error != nil {
		return response.Message(c, fiber.StatusBadRequest, false, "ไม่พบผู้ใช้งานนี้!")
	}

	// Verify OTP
	if user.OTP != input.OTP || time.Now().After(user.OTPExpiry) {
		return response.Message(c, fiber.StatusBadRequest, false, "รหัส OTP หมดอายุแล้ว")
	}

	return response.Message(c, fiber.StatusOK, true, "ยืนยัน OPT สำเร็จ")
}

func (a *AuthController) ResetPassword(c *fiber.Ctx) error {
	var input struct {
		Email       string `json:"email"`
		OTP         string `json:"otp"`
		NewPassword string `json:"new_password"`
	}
	if err := c.BodyParser(&input); err != nil {
		return err
	}

	// Check if email exists
	var user models.User
	if result := a.DB.Where("email = ?", input.Email).First(&user); result.Error != nil {
		return response.Message(c, fiber.StatusBadRequest, false, "ไม่พบผู้ใช้งานนี้!")
	}

	// Verify OTP
	if user.OTP != input.OTP || time.Now().After(user.OTPExpiry) {
		return response.Message(c, fiber.StatusBadRequest, false, "รหัส OTP หมดอายุแล้ว")
	}

	// Update password
	user.Password = utils.GeneratePassword(input.NewPassword)
	a.DB.Save(&user)

	return response.Message(c, fiber.StatusOK, true, "เปลี่ยนรหัสสำเร็จ")
}
