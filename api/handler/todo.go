package handler

import (
	"fmt"
	"net/http"
	"sync"
	"testfiber/api/payload"
	"testfiber/pkg/entities"
	"testfiber/pkg/todo"
	"testfiber/pkg/utility"

	"github.com/gofiber/fiber/v2"
)

func AddTodo(service todo.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var todo entities.Todo
		t := utility.GetTime()

		if err := c.BodyParser(&todo); err != nil {
			c.Status(http.StatusBadRequest)
			c.JSON(payload.ErrorResponse(http.StatusBadRequest, err))
		}

		if err := utility.Check(todo); err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(payload.ErrorResponse(http.StatusBadRequest, err))
		}

		todo.ID = service.ID.Generate()
		todo.CreateAt = t
		todo.UpdateAt = t
		todo.IsActive = true

		if todo.Priority == "" {
			todo.Priority = "very-high"
		}

		var wg = new(sync.WaitGroup)
		wg.Add(2)

		go func() {
			service.Repo.Create(c.Context(), todo)
			wg.Done()
		}()
		go func() {
			service.Sess.Set(c.Context(), fmt.Sprint(todo.ID), todo)
			wg.Done()
		}()

		wg.Wait()
		return c.Status(http.StatusCreated).JSON(payload.SuccessResponse(todo.Map()))
	}
}

func EditTodo(service todo.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var todo entities.Todo
		var cache = make(chan entities.Todo)
		var store = make(chan entities.Todo)
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
		case todo = <-store:
			break
		case todo = <-cache:
			break
		}

		if err := c.BodyParser(&todo); err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(payload.ErrorResponse(http.StatusBadRequest, err))
		}

		todo.UpdateAt = utility.GetTime()
		var wg = new(sync.WaitGroup)
		wg.Add(2)

		go func() {
			service.Repo.Update(c.Context(), todo)
			wg.Done()
		}()

		go func() {
			service.Sess.Set(c.Context(), fmt.Sprint(todo.ID), todo)
			wg.Done()
		}()

		wg.Wait()
		return c.Status(http.StatusOK).JSON(payload.SuccessResponse(todo.Map()))
	}
}

func DeleteTodo(service todo.Service) fiber.Handler {
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

func GetTodo(service todo.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var todo entities.Todo
		var store = make(chan entities.Todo)
		var cache = make(chan entities.Todo)
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
		case todo = <-store:
			break
		case todo = <-cache:
			break
		}

		return c.Status(http.StatusOK).JSON(payload.SuccessResponse(todo.Map()))
	}
}

func GetTodos(service todo.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var todos []fiber.Map
		where := map[string]string{}
		where["activity_group_id"] = c.Query("activity_group_id")

		rows, err := service.Repo.Reads(c.Context(), where)
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
