package payload

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func SuccessResponse(data *fiber.Map) *fiber.Map {
	return &fiber.Map{
		"status":  "Success",
		"message": "Success",
		"data":    data,
	}
}

func SliceSuccessResponse(data []fiber.Map) *fiber.Map {
	if len(data) == 0 {
		return &fiber.Map{
			"status":  "Success",
			"message": "Success",
			"data":    []fiber.Map{},
		}
	}

	return &fiber.Map{
		"status":  "Success",
		"message": "Success",
		"data":    data,
	}
}

func ErrorResponse(code int, err error) *fiber.Map {
	return &fiber.Map{
		"status":  http.StatusText(code),
		"message": err.Error(),
		"data":    nil,
	}
}

func SliceErrorResponse(code int, err error) *fiber.Map {
	return &fiber.Map{
		"status":  http.StatusText(code),
		"message": err.Error(),
		"data":    fiber.Map{},
	}
}
