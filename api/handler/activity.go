package handler

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"testfiber/api/payload"
	"testfiber/storage/activity"
	"testfiber/storage/entities"
	"testfiber/utility"

	"github.com/gofiber/fiber/v2"
)

func AddActivity(service activity.Service) fiber.Handler {
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
		//set cache
		service.Sess.Set(c.Context(), fmt.Sprint(activity.ID), activity)

		c.Status(http.StatusCreated)
		return c.JSON(payload.SuccessResponse(activity.Map()))
	}
}

func EditActivity(service activity.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var activity entities.Activity
		var cache = true
		var time = utility.GetTime()

		// get cache
		err := service.Sess.Get(c.Context(), c.Params("id"), &activity)
		if err != nil {
			cache = false
			log.Println(err)
		}

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
		if err := service.Repo.Update(&activity, time); err != nil {
			c.Status(http.StatusNotFound)
			return c.JSON(payload.ErrorResponse(http.StatusNotFound, err))
		}

		activity.UpdateAt = time

		if !cache {
			where := map[string]string{"id": c.Params("id")}

			//get activity from database
			res, err := service.Repo.Read(where)
			if err != nil {
				c.Status(http.StatusNotFound)
				return c.JSON(payload.ErrorResponse(http.StatusNotFound, err))
			}

			service.Sess.Set(c.Context(), fmt.Sprint(activity.ID), *res)

			c.Status(http.StatusOK)
			return c.JSON(payload.SuccessResponse(res.Map()))
		}

		service.Sess.Set(c.Context(), fmt.Sprint(activity.ID), activity)

		c.Status(http.StatusOK)
		return c.JSON(payload.SuccessResponse(activity.Map()))

	}
}

func DeleteActivity(service activity.Service) fiber.Handler {
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

func GetActivity(service activity.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var activity entities.Activity
		var cache = true
		err := service.Sess.Get(c.Context(), c.Params("id"), &activity)
		if err != nil {
			cache = false
			log.Println(err)
		}

		if !cache {

			where := map[string]string{"id": c.Params("id")}

			//get activity from database
			res, err := service.Repo.Read(where)
			if err != nil {
				c.Status(http.StatusNotFound)
				return c.JSON(payload.ErrorResponse(http.StatusNotFound, err))
			}

			c.Status(http.StatusOK)
			return c.JSON(payload.SuccessResponse(res.Map()))

		}

		c.Status(http.StatusOK)
		return c.JSON(payload.SuccessResponse(activity.Map()))
	}
}

func GetActivities(service activity.Service) fiber.Handler {
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
		return c.JSON(payload.SliceSuccessResponse(activities))
	}
}
