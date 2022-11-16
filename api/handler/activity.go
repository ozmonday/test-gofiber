package handler

import (
	"errors"
	"net/http"
	"strconv"
	"testfiber/api/payload"
	"testfiber/storage/activitiy"
	"testfiber/storage/entities"

	"github.com/gofiber/fiber/v2"
)

func AddActivity(service activitiy.Service) fiber.Handler {
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

		if err := service.Repo.Create(&activity); err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(payload.ErrorResponse(http.StatusInternalServerError, err))
		}

		c.Status(http.StatusCreated)
		return c.JSON(payload.SuccessResponse(activity.Map()))
	}
}

func EditActivity(service activitiy.Service) fiber.Handler {
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
		if err := service.Repo.Update(&activity); err != nil {
			c.Status(http.StatusNotFound)
			return c.JSON(payload.ErrorResponse(http.StatusNotFound, err))
		}

		c.Status(http.StatusOK)
		return c.JSON(payload.SuccessResponse(activity.Map()))

	}
}

func DeleteActivity(service activitiy.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		where := map[string]string{"id": c.Params("id")}

		// delete Activity from database
		if err := service.Repo.Delete(where); err != nil {
			c.Status(http.StatusNotFound)
			return c.JSON(payload.ErrorResponse(http.StatusNotFound, err))
		}

		c.Status(http.StatusOK)
		return c.JSON(payload.SuccessResponse(&fiber.Map{}))
	}
}

func GetActivity(service activitiy.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		where := map[string]string{"id": c.Params("id")}

		//get activity from database
		activity, err := service.Repo.Read(where)
		if err != nil {
			c.Status(http.StatusNotFound)
			return c.JSON(payload.ErrorResponse(http.StatusNotFound, err))
		}

		c.Status(http.StatusOK)
		return c.JSON(payload.SuccessResponse(activity.Map()))
	}

}

func GetActivities(service activitiy.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var activities []fiber.Map
		where := map[string]string{}

		// get all data from
		rows, err := service.Repo.Reads(where)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(payload.ErrorResponse(http.StatusInternalServerError, err))
		}

		for _, v := range *rows {
			activities = append(activities, *v.Map())
		}

		c.Status(http.StatusOK)
		return c.JSON(payload.SliceSuccessResponse(&activities))
	}
}
