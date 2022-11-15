package handler

import (
	"errors"
	"net/http"
	"strconv"
	"testfiber/api/payload"
	"testfiber/storage/entities"

	"github.com/gofiber/fiber/v2"
)

func AddActivity() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var activity entities.Activity
		if err := c.BodyParser(&activity); err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(payload.ErrorResponse(http.StatusBadRequest, err))
		}

		if activity.Title == "" {
			c.Status(http.StatusBadRequest)
			return c.JSON(payload.ErrorResponse(http.StatusBadRequest, errors.New("title cannot be null")))
		}

		if activity.Email == "" {
			c.Status(http.StatusBadRequest)
			return c.JSON(payload.ErrorResponse(http.StatusBadRequest, errors.New("email cannot be null")))
		}
		//add activity to database

		c.Status(http.StatusCreated)
		return c.JSON(payload.SuccessResponse(activity.Map()))
	}
}

func EditActivity() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var activity entities.Activity
		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(payload.ErrorResponse(http.StatusBadRequest, err))
		}

		activity.ID = int64(id)

		if err := c.BodyParser(&activity); err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(payload.ErrorResponse(http.StatusBadRequest, err))
		}

		// update activity from database

		c.Status(http.StatusOK)
		return c.JSON(payload.SuccessResponse(activity.Map()))

	}
}

func DeleteActivity() fiber.Handler {
	return func(c *fiber.Ctx) error {
		_, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(payload.ErrorResponse(http.StatusBadRequest, err))
		}

		// delete Activity from database

		c.Status(http.StatusOK)
		return c.JSON(payload.SuccessResponse(&fiber.Map{}))
	}
}

func GetActivity() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var activity *entities.Activity
		_, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(payload.ErrorResponse(http.StatusBadRequest, err))
		}

		//get activity from database

		c.Status(http.StatusOK)
		return c.JSON(payload.SuccessResponse(activity.Map()))
	}

}

func GetActivities() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var activities []fiber.Map

		// get all data from

		c.Status(http.StatusOK)
		return c.JSON(payload.SliceSuccessResponse(&activities))
	}
}
