package entities

import (
	"github.com/gofiber/fiber/v2"
)

type Activity struct {
	ID       int64  `json:"id" `
	Title    string `json:"title" behav:"required"`
	Email    string `json:"email" behav:"required"`
	CreateAt string `json:"created_at"`
	UpdateAt string `json:"updated_at"`
	DeleteAt string `json:"deleted_at"`
}

func (a *Activity) Map() *fiber.Map {
	result := fiber.Map{
		"id":         a.ID,
		"title":      a.Title,
		"email":      a.Email,
		"created_at": a.CreateAt,
		"updated_at": a.UpdateAt,
		"deleted_at": a.DeleteAt,
	}

	if a.DeleteAt == "" {
		result["deleted_at"] = nil
	}
	return &result
}
