package dto

import (
	"github.com/gofiber/fiber/v2"
)

const (
	FieldBadFormat     = "FIELD_BADFORMAT"
	FieldIncorrect     = "FIELD_INCORRECT"
	ServiceUnavailable = "SERVICE_UNAVAILABLE"
	InternalError      = "Service is currently unavailable. Please try again later."
	NoContent          = "Not Data"
	ContentError       = "Service is available, but no data with this ID"
)

type Response struct {
	Status string `json:"status"`
	Error  *Error `json:"error,omitempty"`
	Data   any    `json:"data,omitempty"`
}

type Error struct {
	Code string `json:"code"`
	Desc string `json:"desc"`
}

func BadResponseError(ctx *fiber.Ctx, code, desc string) error {
	return ctx.Status(fiber.StatusBadRequest).JSON(Response{
		Status: "error",
		Error: &Error{
			Code: code,
			Desc: desc,
		},
	})
}

func InternalServerError(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusInternalServerError).JSON(Response{
		Status: "error",
		Error: &Error{
			Code: ServiceUnavailable,
			Desc: InternalError,
		},
	})
}

func NotFound(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusNotFound).JSON(Response{
		Status: "error",
		Error: &Error{
			Code: NoContent,
			Desc: ContentError,
		},
	})
}
