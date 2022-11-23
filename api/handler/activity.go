package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"testfiber/api/payload"
	"testfiber/pkg/activity"
	"testfiber/pkg/entities"
	"testfiber/pkg/utility"

	"github.com/gofiber/fiber/v2"
)

func AddActivity(service activity.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var activity entities.Activity
		//err := make(chan error)
		if err := c.BodyParser(&activity); err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(payload.ErrorResponse(http.StatusBadRequest, err))
		}

		if err := utility.Check(activity); err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(payload.ErrorResponse(http.StatusBadRequest, err))
		}

		if err := service.Repo.Create(&activity); err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(payload.ErrorResponse(http.StatusInternalServerError, err))
		}

		//set cache
		go service.Sess.Set(c.Context(), fmt.Sprint(activity.ID), activity)

		c.Status(http.StatusCreated)
		return c.JSON(payload.SuccessResponse(activity.Map()))
	}
}

func EditActivity(service activity.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var activity entities.Activity
		var cache bool
		var time = utility.GetTime()

		// get cache
		if err := service.Sess.Get(c.Context(), c.Params("id"), &activity); err == nil {
			cache = true
		}

		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(payload.ErrorResponse(http.StatusBadRequest, err))
		}

		if err := c.BodyParser(&activity); err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(payload.ErrorResponse(http.StatusBadRequest, err))
		}

		activity.ID = int64(id)
		activity.UpdateAt = time

		// update activity from database
		if err := service.Repo.Update(&activity, time); err != nil {
			c.Status(http.StatusNotFound)
			return c.JSON(payload.ErrorResponse(http.StatusNotFound, err))
		}

		if cache {
			go service.Sess.Set(c.Context(), fmt.Sprint(activity.ID), activity)

			c.Status(http.StatusOK)
			return c.JSON(payload.SuccessResponse(activity.Map()))
		}

		where := map[string]string{"id": c.Params("id")}
		res, err := service.Repo.Read(where)
		if err != nil {
			c.Status(http.StatusNotFound)
			return c.JSON(payload.ErrorResponse(http.StatusNotFound, err))
		}

		go service.Sess.Set(c.Context(), fmt.Sprint(activity.ID), *res)

		c.Status(http.StatusOK)
		return c.JSON(payload.SuccessResponse(res.Map()))
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
		if err := service.Sess.Get(c.Context(), c.Params("id"), &activity); err == nil {
			c.Status(http.StatusOK)
			return c.JSON(payload.SuccessResponse(activity.Map()))
		}

		where := map[string]string{"id": c.Params("id")}
		res, err := service.Repo.Read(where)
		if err != nil {
			c.Status(http.StatusNotFound)
			return c.JSON(payload.ErrorResponse(http.StatusNotFound, err))
		}

		c.Status(http.StatusOK)
		return c.JSON(payload.SuccessResponse(res.Map()))

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
