package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"testfiber/api/payload"
	"testfiber/pkg/entities"
	"testfiber/pkg/todo"
	"testfiber/pkg/utility"

	"github.com/gofiber/fiber/v2"
)

func AddTodo(service todo.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var todo entities.Todo
		if err := c.BodyParser(&todo); err != nil {
			c.Status(http.StatusBadRequest)
			c.JSON(payload.ErrorResponse(http.StatusBadRequest, err))
		}

		if err := utility.Check(todo); err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(payload.ErrorResponse(http.StatusBadRequest, err))
		}

		if err := service.Repo.Create(&todo); err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(payload.ErrorResponse(http.StatusInternalServerError, err))
		}
		//set cache
		go service.Sess.Set(c.Context(), fmt.Sprint(todo.ID), todo)

		c.Status(http.StatusCreated)
		return c.JSON(payload.SuccessResponse(todo.Map()))
	}
}

func EditTodo(service todo.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var todo entities.Todo
		var cache bool
		var time = utility.GetTime()

		if err := service.Sess.Get(c.Context(), c.Params("id"), &todo); err == nil {
			cache = true
		}

		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(payload.ErrorResponse(http.StatusBadRequest, err))
		}

		if err := c.BodyParser(&todo); err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(payload.ErrorResponse(http.StatusBadRequest, err))
		}

		todo.ID = int64(id)
		todo.UpdateAt = time

		if err := service.Repo.Update(&todo, time); err != nil {
			c.Status(http.StatusNotFound)
			return c.JSON(payload.ErrorResponse(http.StatusNotFound, err))
		}

		if cache {
			go service.Sess.Set(c.Context(), fmt.Sprint(todo.ID), todo)

			c.Status(http.StatusOK)
			return c.JSON(payload.SuccessResponse(todo.Map()))
		}

		where := map[string]string{"id": c.Params("id")}

		res, err := service.Repo.Read(where)
		if err != nil {
			c.Status(http.StatusNotFound)
			return c.JSON(payload.ErrorResponse(http.StatusNotFound, err))
		}

		go service.Sess.Set(c.Context(), fmt.Sprint(todo.ID), *res)

		c.Status(http.StatusOK)
		return c.JSON(payload.SuccessResponse(res.Map()))
	}
}

func DeleteTodo(service todo.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		where := map[string]string{"id": c.Params("id")}

		if err := service.Repo.Delete(where); err != nil {
			c.Status(http.StatusNotFound)
			return c.JSON(payload.ErrorResponse(http.StatusNotFound, err))
		}

		c.Status(http.StatusOK)
		return c.JSON(payload.SuccessResponse(&fiber.Map{}))
	}
}

func GetTodo(service todo.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var todo entities.Todo

		if err := service.Sess.Get(c.Context(), c.Params("id"), &todo); err == nil {
			c.Status(http.StatusOK)
			return c.JSON(payload.SuccessResponse(todo.Map()))
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

func GetTodos(service todo.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var todos []fiber.Map
		where := map[string]string{}
		where["activity_group_id"] = c.Query("activity_group_id")

		rows, err := service.Repo.Reads(where)
		if err != nil {
			c.Status(http.StatusNotFound)
			return c.JSON(payload.SliceErrorResponse(http.StatusNotFound, err))
		}

		for _, v := range *rows {
			todos = append(todos, *v.Map())
		}

		c.Status(http.StatusOK)
		return c.JSON(payload.SliceSuccessResponse(todos))
	}
}
