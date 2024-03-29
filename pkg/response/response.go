package response

import (
	"math"

	"github.com/gofiber/fiber/v2"
)

type ResponseMessage struct {
	StatusCode int    `json:"statusCode"`
	TaskStatus bool   `json:"taskStatus"`
	Message    string `json:"message"`
}

type ResponseData struct {
	StatusCode int         `json:"statusCode"`
	TaskStatus bool        `json:"taskStatus"`
	Data       interface{} `json:"data"`
	Pagin      *Pagination `json:"pagin,omitempty"`
}

type Pagination struct {
	PageNumber  int `json:"pageNumber"`
	PageSize    int `json:"pageSize"`
	TotalPages  int `json:"totalPages"`
	TotalRecord int `json:"totalRecord"`
}

func Message(ctx *fiber.Ctx, statusCode int, taskStatus bool, message string) error {
	response := ResponseMessage{
		StatusCode: statusCode,
		TaskStatus: taskStatus,
		Message:    message,
	}
	return ctx.Status(statusCode).JSON(response)
}

func SendData(ctx *fiber.Ctx, statusCode int, taskStatus bool, data interface{}, pagin *Pagination) error {
	response := ResponseData{
		StatusCode: statusCode,
		TaskStatus: taskStatus,
		Data:       data,
		Pagin:      pagin,
	}
	if pagin != nil {
		pagin.TotalPages = calculateTotalPages(pagin.TotalRecord, pagin.PageSize)
		response.Pagin = pagin
	}
	return ctx.Status(statusCode).JSON(response)
}

func calculateTotalPages(totalRecord, pageSize int) int {
	return int(math.Ceil(float64(totalRecord) / float64(pageSize)))
}
