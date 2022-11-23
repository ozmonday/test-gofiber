package handler

import (
	"fmt"
	"net/http"
	"sync"
	"testfiber/api/payload"
	"testfiber/pkg/activity"
	"testfiber/pkg/entities"
	"testfiber/pkg/utility"

	"github.com/gofiber/fiber/v2"
)

func AddActivity(service activity.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var activity entities.Activity
		t := utility.GetTime()

		if err := c.BodyParser(&activity); err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(payload.ErrorResponse(http.StatusBadRequest, err))
		}

		if err := utility.Check(activity); err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(payload.ErrorResponse(http.StatusBadRequest, err))
		}
		activity.ID = service.ID.Generate()
		activity.CreateAt = t
		activity.UpdateAt = t

		var wg = new(sync.WaitGroup)
		wg.Add(2)
		go func() {
			service.Repo.Create(c.Context(), activity)
			wg.Done()
		}()

		go func() {
			service.Sess.Set(c.Context(), fmt.Sprint(activity.ID), activity)
			wg.Done()
		}()

		wg.Wait()
		return c.Status(http.StatusCreated).JSON(payload.SuccessResponse(activity.Map()))
	}
}

func EditActivity(service activity.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var activity entities.Activity
		var cache = make(chan entities.Activity)
		var store = make(chan entities.Activity)
		var errc = make(chan error)

		go func() {
			res, err := service.Sess.Get(c.Context(), c.Params("id"))
			if err != nil {
				return
			}
			cache <- *res
			close(cache)
		}()

		go func() {
			where := map[string]string{"id": c.Params("id")}
			res, err := service.Repo.Read(c.Context(), where)
			if err != nil {
				errc <- err
				close(errc)
				return
			}
			store <- *res
			close(store)
		}()

		select {
		case err := <-errc:
			c.Status(http.StatusNotFound)
			return c.JSON(payload.ErrorResponse(http.StatusNotFound, err))
		case activity = <-store:
			break
		case activity = <-cache:
			break
		}

		if err := c.BodyParser(&activity); err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(payload.ErrorResponse(http.StatusBadRequest, err))
		}

		activity.UpdateAt = utility.GetTime()

		var wg = new(sync.WaitGroup)
		wg.Add(2)

		go func() {
			service.Repo.Update(c.Context(), activity)
			wg.Done()
		}()

		go func() {
			service.Sess.Set(c.Context(), fmt.Sprint(activity.ID), activity)
			wg.Done()
		}()

		wg.Wait()
		return c.Status(http.StatusOK).JSON(payload.SuccessResponse(activity.Map()))
	}
}

func DeleteActivity(service activity.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		where := map[string]string{"id": c.Params("id")}

		if err := service.Repo.Delete(c.Context(), where); err != nil {
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
		var store = make(chan entities.Activity)
		var cache = make(chan entities.Activity)
		var errc = make(chan error)

		go func() {
			res, err := service.Sess.Get(c.Context(), c.Params("id"))
			if err != nil {
				return
			}
			cache <- *res
			close(cache)
		}()

		go func() {
			where := map[string]string{"id": c.Params("id")}
			res, err := service.Repo.Read(c.Context(), where)
			if err != nil {
				errc <- err
				close(errc)
				return
			}
			store <- *res
			close(store)
		}()

		select {
		case err := <-errc:
			c.Status(http.StatusNotFound)
			return c.JSON(payload.ErrorResponse(http.StatusNotFound, err))
		case activity = <-store:
			break
		case activity = <-cache:
			break
		}

		return c.Status(http.StatusOK).JSON(payload.SuccessResponse(activity.Map()))
	}
}

func GetActivities(service activity.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var activities []fiber.Map
		where := map[string]string{}

		// get all data from
		rows, err := service.Repo.Reads(c.Context(), where)
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
