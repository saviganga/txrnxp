package utils

import (
	"github.com/gofiber/fiber/v2"
)

func NoDataSuccessResponse(ctx *fiber.Ctx, message string) error {
	resp := make(map[string]interface{})
	resp["status"] = "SUCCESS"
	resp["type"] = "Ok"
	resp["message"] = message
	return ctx.Status(fiber.StatusOK).JSON(resp)
}

func SuccessResponse(ctx *fiber.Ctx, data interface{}, message string) error {
	resp := make(map[string]interface{})
	resp["status"] = "SUCCESS"
	resp["type"] = "Ok"
	resp["message"] = message
	resp["data"] = data
	return ctx.Status(fiber.StatusOK).JSON(resp)
}

func PaginatedSuccessResponse(ctx *fiber.Ctx, data interface{}, message string) error {
	resp := data
	return ctx.Status(fiber.StatusOK).JSON(resp)
}

func CreatedResponse(ctx *fiber.Ctx, data interface{}, message string) error {
	resp := make(map[string]interface{})
	resp["status"] = "SUCCESS"
	resp["type"] = "Created"
	resp["message"] = message
	resp["data"] = data
	return ctx.Status(fiber.StatusCreated).JSON(resp)
}

func BadRequestResponse(ctx *fiber.Ctx, message string) error {
	resp := make(map[string]interface{})
	resp["status"] = "FAILED"
	resp["type"] = "Bad Request"
	resp["message"] = message
	return ctx.Status(fiber.StatusBadRequest).JSON(resp)
}

func UnauthorisedResponse(ctx *fiber.Ctx, message string) error {
	resp := make(map[string]interface{})
	resp["status"] = "FAILED"
	resp["type"] = "Unauthorised"
	resp["message"] = message
	return ctx.Status(fiber.StatusUnauthorized).JSON(resp)
}

func ForbiddenResponse(ctx *fiber.Ctx, message string) error {
	resp := make(map[string]interface{})
	resp["status"] = "FAILED"
	resp["type"] = "Unauthorised"
	resp["message"] = message
	return ctx.Status(fiber.StatusForbidden).JSON(resp)
}

func ServerErrorResponse(ctx *fiber.Ctx) error {
	resp := make(map[string]interface{})
	resp["status"] = "FAILED"
	resp["type"] = "Internal Server Error"
	resp["message"] = "Failure, Invalid Params, Database Error or Server Error"
	return ctx.Status(fiber.StatusInternalServerError).JSON(resp)
}
